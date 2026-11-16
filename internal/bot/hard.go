package bot

import (
	"context"
	"errors"
	"math/rand/v2"
	"runtime"
	"slices"
	"time"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The hard level (ADR 030, Decision 9): many imagined worlds; each
// candidate is played out on copies of them to the end of the turn,
// every seat by the policy, and scored by settledScore with the
// opponents weighted by strength (Decision 7.6). Rounds of playouts
// halve the candidates; the tries go to the best ones. No mistakes.
const (
	// playoutSteps caps a playout (Decision 9.3).
	playoutSteps = 120
	// defaultThink is the time a budget with neither cap thinks
	// (Decision 9.8).
	defaultThink = 1500 * time.Millisecond
	// yieldSteps: a playout yields the core every this many steps, as
	// well as after it ends; a long playout would else hold it.
	yieldSteps = 16
	// playoutStream is the PCG stream of a playout's random source.
	playoutStream = 0x91a7
)

// hardPick is the hard level's move for in.Seat: w is World(0) and
// allowed its Allowed. It stops at the budget or when ctx ends, and
// then plays the best so far (Decision 9.9); with nothing tried yet, a
// pass or an end of turn when there is one, else the first candidate.
func hardPick(ctx context.Context, in Input, w *engine.Game, allowed []engine.Intent) (engine.Intent, error) {
	p := in.Seat
	cs := candidates(w, p, keepPlays(w, p, allowed), in.Rand)
	// A pass or an end of turn first: doing nothing is the baseline a
	// move must beat (Decision 8.6), and it stays to the end.
	slices.SortStableFunc(cs, func(a, b engine.Intent) int { return quietRank(a) - quietRank(b) })
	if len(cs) == 1 {
		return cs[0], nil
	}
	s := search{
		done: ctx.Done(), in: in, p: p, cs: cs,
		sc:    scoring{lead: true, answers: answers(w, p, allowed)},
		tries: make([][]float64, len(cs)),
		cap:   in.Budget.Playouts,
		start: time.Now(),
	}
	if s.cap == 0 && in.Budget.Time == 0 {
		s.until = s.start.Add(defaultThink)
	} else if in.Budget.Time > 0 {
		s.until = s.start.Add(in.Budget.Time)
	}
	s.noise = noiseOpen
	if len(w.Stack) > 0 {
		s.noise = noiseStack
	}
	if err := s.run(w); err != nil {
		return engine.Intent{}, err
	}
	return cs[s.best()], nil
}

// search is one hard move's state.
type search struct {
	// done is the move's ctx.Done().
	done  <-chan struct{}
	in    Input
	p     engine.PlayerID
	cs    []engine.Intent
	sc    scoring
	noise float64
	// tries[c] are candidate c's playout scores, one per world it met;
	// every candidate still in meets the same worlds, in order.
	tries [][]float64
	// alive are the candidates still in, in the order of cs.
	alive []int
	// cap is the most playouts (0: none); until the time to stop (zero:
	// none).
	cap, played int
	start       time.Time
	until       time.Time
	// worlds counts the worlds asked for.
	worlds int
}

// errStop: the budget or ctx ended the search.
var errStop = errors.New("bot: search stopped")

// run plays the rounds. Each round gives every candidate still in an
// equal share of the budget, world by world; then the better half of
// the moves stays, by judge. Doing nothing (cs[0] when quiet) is the
// baseline and always stays; the last round plays it against the one
// move left.
func (s *search) run(w0 *engine.Game) error {
	s.alive = make([]int, len(s.cs))
	for i := range s.alive {
		s.alive[i] = i
	}
	rounds := 1
	for n := 1; n < s.moves(); n *= 2 {
		rounds++
	}
	for round := 0; ; round++ {
		left := rounds - round
		per := s.perRound(left)
		for k := 0; k < per; k++ {
			err := s.world(w0)
			if errors.Is(err, errStop) {
				return nil
			}
			if err != nil {
				return err
			}
			if k == 0 && round == 0 && !s.until.IsZero() {
				// A time budget: the first world tells how long a
				// playout takes.
				per = 1 + s.perRound(left)
			}
		}
		if s.moves() <= 1 {
			return nil
		}
		s.halve()
	}
}

// quiet: cs[0] is doing nothing, the baseline of judge.
func (s *search) quiet() bool { return quietRank(s.cs[0]) == 0 }

// moves counts the candidates still in besides doing nothing.
func (s *search) moves() int {
	if s.quiet() && slices.Contains(s.alive, 0) {
		return len(s.alive) - 1
	}
	return len(s.alive)
}

// perRound is how many worlds each candidate still in meets this
// round, with left rounds to go: an even share of the playouts left.
// With a time budget the playouts left come from the time left and the
// time a playout took so far, at most the cap; before any playout, one
// world.
func (s *search) perRound(left int) int {
	budget := s.cap - s.played
	if !s.until.IsZero() {
		if s.played == 0 {
			return 1
		}
		each := time.Since(s.start) / time.Duration(s.played)
		byTime := int(time.Until(s.until) / max(each, 1))
		if s.cap == 0 || byTime < budget {
			budget = byTime
		}
	}
	return max(1, budget/(len(s.alive)*max(left, 1)))
}

// world plays every candidate still in once on the next world. It
// returns errStop when the budget or ctx ends first.
func (s *search) world(w0 *engine.Game) error {
	w := w0
	if s.worlds > 0 {
		if s.stop() {
			return errStop
		}
		var err error
		if w, err = s.in.World(s.worlds); err != nil {
			return err
		}
		runtime.Gosched()
	}
	s.worlds++
	sc := s.sc
	sc.dig = digValue(w)
	seed := s.in.Rand.Uint64()
	for _, c := range s.alive {
		if s.stop() {
			return errStop
		}
		score, err := playout(w, s.p, s.cs[c], seed, sc)
		switch {
		case errors.Is(err, errRefused):
			// Every candidate is allowed; a refused one is a bug,
			// never played.
			score = -2 * scoreWin
		case err != nil:
			return err
		}
		s.tries[c] = append(s.tries[c], score)
		s.played++
		// Let the room and other games run: one core may hold many
		// tables (Decision 17).
		runtime.Gosched()
	}
	return nil
}

// stop: the playouts or the time are used up, or ctx ended.
func (s *search) stop() bool {
	select {
	case <-s.done:
		return true
	default:
	}
	return s.cap > 0 && s.played >= s.cap || !s.until.IsZero() && !time.Now().Before(s.until)
}

// halve keeps the better half of the moves still in, by judge on the
// worlds they all met; doing nothing always stays.
func (s *search) halve() {
	scores := s.judged()
	var order []int
	for _, c := range s.alive {
		if c != 0 || !s.quiet() {
			order = append(order, c)
		}
	}
	slices.SortStableFunc(order, func(a, b int) int { return rank(scores[b] - scores[a]) })
	keep := order[:(len(order)+1)/2]
	s.alive = s.alive[:0]
	for c := range s.cs {
		if c == 0 && s.quiet() || slices.Contains(keep, c) {
			s.alive = append(s.alive, c)
		}
	}
}

// best is the candidate to play: the top of judge among those still
// in, on equal scores the first (a pass or an end of turn). With
// nothing tried yet, the first candidate.
func (s *search) best() int {
	scores := s.judged()
	best := -1
	for _, c := range s.alive {
		if _, ok := scores[c]; ok && (best < 0 || scores[c]-scores[best] > scoreEps) {
			best = c
		}
	}
	return max(best, 0)
}

// judged is judge's score of each candidate still in that was tried,
// on the worlds all of them met: a search stopped within a world cuts
// the others to the same worlds.
func (s *search) judged() map[int]float64 {
	var cs []engine.Intent
	var tries [][]float64
	var idx []int
	n := -1
	for _, c := range s.alive {
		if len(s.tries[c]) > 0 && (n < 0 || len(s.tries[c]) < n) {
			n = len(s.tries[c])
		}
	}
	for _, c := range s.alive {
		if len(s.tries[c]) == 0 {
			continue
		}
		cs = append(cs, s.cs[c])
		tries = append(tries, s.tries[c][:n])
		idx = append(idx, c)
	}
	out := make(map[int]float64, len(idx))
	for i, v := range judge(cs, tries, s.noise, noiseClip) {
		out[idx[i]] = v
	}
	return out
}

// rank orders by a score difference d: above scoreEps first.
func rank(d float64) int {
	switch {
	case d > scoreEps:
		return 1
	case d < -scoreEps:
		return -1
	}
	return 0
}

// playout is the score of candidate in on world w: applied to a copy,
// played on to the end of the turn by the policy for every seat (a
// choice of seat p gets the policy's best option, as in a settle), at
// most playoutSteps, and scored by settledScore (Decision 9.2 – 9.4).
// seed seeds the policy's choices; the dice are the world's own.
func playout(w *engine.Game, p engine.PlayerID, in engine.Intent, seed uint64, sc scoring) (float64, error) {
	g, err := fork(w)
	if err != nil {
		return 0, err
	}
	ev, err := g.Apply(in)
	if err != nil {
		return 0, errors.Join(errRefused, err)
	}
	var o outcome
	o.see(ev, p)
	var src rand.PCG
	src.Seed(seed, playoutStream)
	rng := rand.New(&src) //nolint:gosec // a playout's choices, not security
	turn := g.Turn.Number
	for step := range playoutSteps {
		if g.Over || g.Turn.Number != turn {
			break
		}
		if step%yieldSteps == yieldSteps-1 {
			runtime.Gosched()
		}
		pr := g.Prompt()
		who := pr.Who()
		if len(who) == 0 {
			break
		}
		q := who[0]
		allowed := g.Allowed(q)
		if len(allowed) == 0 {
			break
		}
		var next engine.Intent
		if q == p && pr.Kind == engine.PromptChoose {
			next = policyBest(g, q, allowed, rng)
		} else {
			next = policyPick(g, q, allowed, rng)
		}
		evs, ok := playOn(g, next)
		if !ok {
			break
		}
		o.see(evs, p)
	}
	return settledScore(g, p, o, sc), nil
}
