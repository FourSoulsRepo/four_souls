package bot

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The standard level (ADR 030, Decision 8): one imagined world; each
// candidate is tried on a copy, settled and scored; the best is played,
// with one mistake in ten.
const (
	// settleSteps caps a settle (Decision 8.4). A fight is many rolls,
	// each a stack item every seat passes: a cap of 40 cut off a fifth of
	// 2-player fights and half of 4-player ones; 150 cuts off almost
	// none.
	settleSteps = 150
	// diceDraws is how often a candidate whose settle rolled dice is
	// settled (Decision 8.5): the world's own dice, then the copy's dice
	// reseeded, the same seeds for every candidate. One settle judges a
	// fight by luck.
	diceDraws = 6
	// mistakeOdds: one move in this many is picked by chance among the
	// best mistakeTop, weighted mistakeTop, …, 1 by rank (Decision 8.9).
	mistakeOdds = 10
	mistakeTop  = 3
	// scoreEps: scores this close are equal; then the earlier candidate
	// wins, and a pass or an end of turn comes first.
	scoreEps = 1e-9
	// settleStream is the PCG stream of a settle's random source.
	settleStream = 0x5e771e
)

// errNoPass: the prompt is not one to pass.
var errNoPass = errors.New("bot: not a priority prompt")

// errRefused: the engine refused a candidate.
var errRefused = errors.New("bot: candidate refused")

// standardPick is the standard level's move for seat p on world w: the
// best score; one move in mistakes (0: never) by chance among the best
// mistakeTop.
func standardPick(w *engine.Game, p engine.PlayerID, allowed []engine.Intent, rng *rand.Rand, mistakes int) (engine.Intent, error) {
	cs := candidates(w, p, keepPlays(w, p, allowed), rng)
	// A pass or an end of turn first: on equal scores, doing nothing
	// beats a move that changes nothing.
	slices.SortStableFunc(cs, func(a, b engine.Intent) int { return quietRank(a) - quietRank(b) })
	draws := make([]uint64, diceDraws)
	for i := range draws {
		draws[i] = rng.Uint64()
	}
	dig := digValue(w)
	tries := make([][]float64, len(cs))
	for i, in := range cs {
		t, err := try(w, p, in, draws, dig)
		switch {
		case errors.Is(err, errRefused):
			// Every candidate is allowed; a refused one is a bug, never
			// played.
			t = []float64{-2 * scoreWin}
		case err != nil:
			return engine.Intent{}, err
		}
		tries[i] = t
	}
	noise := noiseOpen
	if len(w.Stack) > 0 {
		noise = noiseStack
	}
	scores := judge(cs, tries, noise, 0)
	order := make([]int, len(cs))
	for i := range order {
		order[i] = i
	}
	slices.SortStableFunc(order, func(a, b int) int {
		switch d := scores[b] - scores[a]; {
		case d > scoreEps:
			return 1
		case d < -scoreEps:
			return -1
		}
		return 0
	})
	best := order[0]
	if mistakes > 0 && len(order) > 1 && rng.IntN(mistakes) == 0 {
		top := min(len(order), mistakeTop)
		ws := make([]float64, top)
		for i := range ws {
			ws[i] = float64(top - i)
		}
		best = order[weighted(ws, rng)]
	}
	return cs[best], nil
}

// judge is each candidate's score from its draws: their mean. When
// there is a way to do nothing (a pass, or an end of turn; first in cs
// after the sort), a move must beat it (Decision 8.6): the move's gain
// over doing nothing, draw by draw on the same seeds, counts less wAct
// and less noise times its standard error. So a move whose gain is
// noise is not made. clip, when not 0, bounds each gain in the
// standard error to ±clip (the mean keeps it whole): a game won or lost
// in one draw of many then does not hide what the others show.
func judge(cs []engine.Intent, tries [][]float64, noise, clip float64) []float64 {
	scores := make([]float64, len(cs))
	for i, t := range tries {
		scores[i] = mean(t)
	}
	if len(cs) == 0 || quietRank(cs[0]) != 0 {
		return scores
	}
	nothing := tries[0]
	gain := make([]float64, 0, diceDraws)
	for i := range cs {
		if quietRank(cs[i]) == 0 {
			continue
		}
		gain = gain[:0]
		for d, s := range tries[i] {
			gain = append(gain, s-nothing[min(d, len(nothing)-1)])
		}
		m := mean(gain)
		if clip > 0 {
			for d, x := range gain {
				gain[d] = min(max(x, -clip), clip)
			}
		}
		scores[i] = scores[0] + m - wAct - noise*stdErr(gain)
	}
	return scores
}

// mean is the mean of xs, which is not empty.
func mean(xs []float64) float64 {
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

// stdErr is the standard error of the mean of xs; 0 for one value.
func stdErr(xs []float64) float64 {
	n := len(xs)
	if n < 2 {
		return 0
	}
	m, ss := mean(xs), 0.0
	for _, x := range xs {
		ss += (x - m) * (x - m)
	}
	return math.Sqrt(ss / float64(n-1) / float64(n))
}

// keepPlays drops the options of a choice that take a play back, as
// the policy does (Decision 6): a play taken back is tried again, for
// ever. With nothing else left, it keeps them all.
func keepPlays(w *engine.Game, p engine.PlayerID, allowed []engine.Intent) []engine.Intent {
	pr := w.PromptFor(p)
	if pr.Kind != engine.PromptChoose {
		return allowed
	}
	ws := choiceWeights(w, p, pr, allowed)
	var out []engine.Intent
	for i, in := range allowed {
		if ws[i] > 0 {
			out = append(out, in)
		}
	}
	if len(out) == 0 {
		return allowed
	}
	return out
}

// quietRank is 0 for a pass or an end of turn, 1 for any other move.
func quietRank(in engine.Intent) int {
	if in.Kind == engine.IntentPass || in.Kind == engine.IntentEndTurn {
		return 0
	}
	return 1
}

// try is the scores of candidate in on world w, one per draw: applied
// to a copy and settled. Draw d settles with the policy's random source
// seeded from draws[d]; draw 0 keeps the world's dice, the others
// reseed the copy's dice from draws[d]. A settle that rolled no dice is
// the only draw. Every candidate meets the same draws, so their scores
// differ by the move, not by luck.
func try(w *engine.Game, p engine.PlayerID, in engine.Intent, draws []uint64, dig float64) ([]float64, error) {
	var src rand.PCG
	prng := rand.New(&src) //nolint:gosec // a playout's choices, not security
	out := make([]float64, 0, len(draws))
	for d, seed := range draws {
		g, err := fork(w)
		if err != nil {
			return nil, err
		}
		if d > 0 {
			g.RNG = engine.NewRNG(seed)
		}
		ev, err := g.Apply(in)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errRefused, err)
		}
		var o outcome
		o.see(ev, p)
		src.Seed(seed, settleStream)
		settle(g, p, prng, &o)
		out = append(out, settledScore(g, p, o, scoring{dig: dig}))
		if !o.rolled {
			break
		}
	}
	return out, nil
}

// fork copies a world: the one place the bot copies a game (Decision
// 4.4, TestBotRules).
func fork(w *engine.Game) (*engine.Game, error) {
	return w.Clone()
}

// settle plays g on until the table is quiet: every seat passes, the
// policy answers the other prompts (Decision 8.3); a choice of seat p
// gets the policy's best option, as p will choose it with care. It stops at the game
// end, when the turn ends, when the active player may act on an empty
// stack with no attack or purchase going on, or after settleSteps
// (Decision 8.4). What it sees goes into o.
func settle(g *engine.Game, p engine.PlayerID, rng *rand.Rand, o *outcome) {
	turn := g.Turn.Number
	for range settleSteps {
		if g.Over || g.Turn.Number != turn || quiet(g) {
			break
		}
		pr := g.Prompt()
		q := pr.Who()[0]
		// A priority prompt is passed: tried at once, as Allowed costs
		// about as much as the pass.
		ev, err := []engine.Event(nil), errNoPass
		if pr.Kind == engine.PromptPriority {
			ev, err = g.Apply(engine.Intent{Player: q, Kind: engine.IntentPass})
		}
		if err != nil {
			allowed := g.Allowed(q)
			if len(allowed) == 0 {
				break
			}
			in := policyPick(g, q, allowed, rng)
			if q == p && pr.Kind == engine.PromptChoose {
				in = policyBest(g, q, allowed, rng)
			}
			if ev, err = g.Apply(in); err != nil {
				break
			}
		}
		o.see(ev, p)
	}
}

// outcome is what a settle saw happen, besides the table it left.
type outcome struct {
	rolled bool // a die was rolled
	kills  int  // monsters whose rewards seat p got
}

// see adds the events of one step for seat p.
func (o *outcome) see(ev []engine.Event, p engine.PlayerID) {
	for _, e := range ev {
		switch e.Kind { //nolint:exhaustive // the other events tell nothing here
		case engine.EvDiceRolled:
			o.rolled = true
		case engine.EvRewardsGained:
			if e.Player == p {
				o.kills++
			}
		}
	}
}

// settledScore is the score of a settled copy from seat p's side: the
// table's score as the level scores it (scoreAs), plus sc.dig for each
// monster p killed in the settle (Decision 7.9). Every level scores a
// settle or a playout by it.
func settledScore(g *engine.Game, p engine.PlayerID, o outcome, sc scoring) float64 {
	return scoreAs(g, p, sc) + sc.dig*float64(o.kills)
}

// quiet: a priority prompt on an empty stack, no attack or purchase
// going on: the active player is free to act.
func quiet(g *engine.Game) bool {
	return g.Prompt().Kind == engine.PromptPriority && len(g.Stack) == 0 && !g.Attack.On && !g.Purchase.On
}
