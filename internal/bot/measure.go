package bot

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The measurement of card values (ADR 033, Decision 20): make bot-values
// (cmd/bot_values) calls PlayGain and ItemPlayout on positions of arena
// games. Both score with the flat values (Decision 20.6), so the
// measure never feeds on itself.
const (
	// valueSteps caps an item's playout: to the end of the seat's next
	// turn, a few turns of a 4-player game.
	valueSteps = 1500
	// valueStream is the PCG stream of an item playout's random source.
	valueStream = 0xa1e5
	// ValueClip bounds each paired difference of the measure to one
	// soul either way (Decision 20.4.9): a game won or lost in one copy
	// and not the other is a million points, which one sample in a
	// hundred would turn into the whole mean.
	ValueClip = wSoul
)

// PlayGain is what seat p gains on world(0) by playing loot card id, or
// by using the abilities of item id, over doing nothing (a pass, or an
// end of turn in its action phase), in settledScore units (Decision
// 20.4.1 – 20.4.4, 20.4.10): each such play or use and doing nothing
// are settled as the standard level settles a candidate, on the same
// draws; the gain is the best one's mean gain, draw by draw, each draw's
// gain clipped to ±ValueClip. ok is false when the seat may not play or
// use the card now, or has no way to do nothing. ctx ends it early.
func PlayGain(ctx context.Context, world World, p engine.PlayerID, id engine.ObjectID, draws []uint64) (gain float64, ok bool, err error) {
	if len(draws) == 0 {
		return 0, false, errors.New("bot: no draws")
	}
	w, err := world(0)
	if err != nil {
		return 0, false, fmt.Errorf("bot: world: %w", err)
	}
	var plays []engine.Intent
	nothing, quietOK := engine.Intent{}, false
	for _, in := range w.Allowed(p) {
		switch {
		case quietRank(in) == 0 && !quietOK:
			nothing, quietOK = in, true
		case (in.Kind == engine.IntentPlayLoot || in.Kind == engine.IntentActivate) && len(in.Objects) > 0 && in.Objects[0] == id:
			plays = append(plays, in)
		}
	}
	if len(plays) == 0 || !quietOK {
		return 0, false, nil
	}
	sc := scoring{dig: digValue(w), flat: true}
	base, _, err := try(ctx, w, p, nothing, draws, sc)
	if err != nil {
		return 0, false, err
	}
	best := 0.0
	for i, in := range plays {
		t, _, err := try(ctx, w, p, in, draws, sc)
		if err != nil {
			return 0, false, err
		}
		n := max(len(t), len(base))
		sum := 0.0
		for d := range n {
			sum += min(max(t[min(d, len(t)-1)]-base[min(d, len(base)-1)], -ValueClip), ValueClip)
		}
		if g := sum / float64(n); i == 0 || g > best {
			best = g
		}
	}
	return best, true, nil
}

// ItemPlayout is the score of world(0) for seat p after a playout to
// the end of p's next turn (Decision 20.4.5 – 20.4.7): every seat moves by the
// policy, p's own choices get the policy's best option, the dice come
// from seed and so do the policy's choices; at most valueSteps. It is
// scored by settledScore with the flat values. Two worlds that differ
// by one item give its value, on the same seeds.
func ItemPlayout(world World, p engine.PlayerID, seed uint64) (float64, error) {
	w, err := world(0)
	if err != nil {
		return 0, fmt.Errorf("bot: world: %w", err)
	}
	g, err := fork(w)
	if err != nil {
		return 0, err
	}
	g.RNG = engine.NewRNG(seed)
	var src rand.PCG
	src.Seed(seed, valueStream)
	rng := rand.New(&src) //nolint:gosec // a playout's choices, not security
	var o outcome
	// next: p's next turn has begun; it ends when the turn number moves
	// on from it.
	start, next := g.Turn.Number, -1
	for range valueSteps {
		if g.Over {
			break
		}
		if next < 0 && g.Turn.Number > start && g.Turn.Active == p {
			next = g.Turn.Number
		}
		if next >= 0 && g.Turn.Number != next {
			break
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
		in := policyPick(g, q, allowed, rng)
		if q == p && pr.Kind == engine.PromptChoose {
			in = policyBest(g, q, allowed, rng)
		}
		ev, ok := playOn(g, in)
		if !ok {
			break
		}
		o.see(ev, p)
	}
	return settledScore(g, p, o, scoring{dig: digValue(w), flat: true}), nil
}
