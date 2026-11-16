package bot

import (
	"math/rand/v2"
	"slices"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The candidate limits of Decision 5.
const (
	// maxDiscards is how many sets of cards a discard tries.
	maxDiscards = 8
	// maxOrders is how many random orders an order tries, besides the
	// order now.
	maxOrders = 5
	// tries bounds the random draws spent on finding new sets or orders.
	tries = 4
)

// candidates lists the moves the seat may try (Decision 5): every
// allowed intent, a discard as up to maxDiscards sets of cards, an order
// as the order now and up to maxOrders random ones. Each is accepted by
// Apply as it is.
func candidates(g *engine.Game, p engine.PlayerID, allowed []engine.Intent, rng *rand.Rand) []engine.Intent {
	var out []engine.Intent
	for _, in := range allowed {
		switch in.Kind { //nolint:exhaustive // the other kinds are complete as allowed
		case engine.IntentDiscard:
			out = append(out, discards(in, g.PromptFor(p).Count, rng)...)
		case engine.IntentOrder:
			out = append(out, orders(in, rng)...)
		default:
			out = append(out, in)
		}
	}
	return out
}

// discards are up to maxDiscards different sets of n cards of the
// template's Objects, picked at random; each set keeps the hand's order.
func discards(tmpl engine.Intent, n int, rng *rand.Rand) []engine.Intent {
	var out []engine.Intent
	var seen [][]engine.ObjectID
	for range maxDiscards * tries {
		if len(out) == maxDiscards {
			break
		}
		set := pickSet(tmpl.Objects, n, rng)
		if slices.ContainsFunc(seen, func(s []engine.ObjectID) bool { return slices.Equal(s, set) }) {
			continue
		}
		seen = append(seen, set)
		out = append(out, engine.Intent{Player: tmpl.Player, Kind: tmpl.Kind, Objects: set})
	}
	return out
}

// pickSet is n of ids at random, in the order of ids.
func pickSet(ids []engine.ObjectID, n int, rng *rand.Rand) []engine.ObjectID {
	n = min(max(n, 0), len(ids))
	idx := rng.Perm(len(ids))[:n]
	slices.Sort(idx)
	set := make([]engine.ObjectID, n)
	for i, j := range idx {
		set[i] = ids[j]
	}
	return set
}

// orders are the order now and up to maxOrders other random orders.
func orders(tmpl engine.Intent, rng *rand.Rand) []engine.Intent {
	out := []engine.Intent{tmpl}
	for range maxOrders * tries {
		if len(out) == maxOrders+1 {
			break
		}
		o := shuffled(tmpl.Objects, rng)
		if slices.ContainsFunc(out, func(in engine.Intent) bool { return slices.Equal(in.Objects, o) }) {
			continue
		}
		out = append(out, engine.Intent{Player: tmpl.Player, Kind: tmpl.Kind, Objects: o})
	}
	return out
}

// shuffled is a random order of ids, in a new slice.
func shuffled(ids []engine.ObjectID, rng *rand.Rand) []engine.ObjectID {
	o := slices.Clone(ids)
	rng.Shuffle(len(o), func(i, j int) { o[i], o[j] = o[j], o[i] })
	return o
}

// complete makes one allowed intent ready to send: a discard gets a
// random set of the asked count, an order a random order.
func complete(g *engine.Game, p engine.PlayerID, in engine.Intent, rng *rand.Rand) engine.Intent {
	switch in.Kind { //nolint:exhaustive // the other kinds are complete as allowed
	case engine.IntentDiscard:
		in.Objects = pickSet(in.Objects, g.PromptFor(p).Count, rng)
	case engine.IntentOrder:
		in.Objects = shuffled(in.Objects, rng)
	}
	return in
}
