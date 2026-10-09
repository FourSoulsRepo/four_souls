package rulesengine

import "fmt"

// ActionKind is a change to the game that cards may rewrite (A-05).
type ActionKind int

// The action kinds; more come with combat and effect blocks.
const (
	ActGainCents ActionKind = iota // R-MECH-36
	ActLoseCents                   // R-MECH-42
	ActLoot                        // R-CARD-06
)

// Action is a pending change. It sits in the queue, may be rewritten by
// replacement effects, and only then happens (R-ABIL-29).
type Action struct {
	Kind   ActionKind `json:"kind"`
	Player PlayerID   `json:"player"`
	Object ObjectID   `json:"object,omitempty"`
	Amount int        `json:"amount"`
	// Applied lists replacements already applied; each applies once
	// (R-ABIL-32).
	Applied []ReplacementRef `json:"applied,omitempty"`
}

// ReplacementRef names one replacement of one object in play.
type ReplacementRef struct {
	Object ObjectID `json:"object"`
	Index  int      `json:"index"`
}

// Replacement is a replacement effect created by a card (R-ABIL-29 to
// R-ABIL-31). When says whether it applies to an action; Do returns what
// happens instead: nothing (prevent), the same action changed, or others.
type Replacement struct {
	Text string
	When func(g *Game, self ObjectID, a Action) bool
	Do   func(g *Game, self ObjectID, a Action) []Action
}

// Events about actions.
const (
	EvReplaced  EventKind = "replaced"
	EvLostCents EventKind = "lost_cents"
)

// enqueue adds actions to the end of the queue.
func (g *Game) enqueue(as ...Action) {
	g.Queue = append(g.Queue, as...)
}

// replacementsFor lists replacements of objects in play that apply to a
// and have not been applied to it yet. Order: by object ID, which is
// stable; the affected player chooses the real order (R-ABIL-33).
func (g *Game) replacementsFor(a Action) []ReplacementRef {
	var out []ReplacementRef
	for i := range g.Objects {
		o := &g.Objects[i]
		if o.Zone.Kind != ZoneInPlay {
			continue
		}
		def, ok := g.cards.find(o.Card)
		if !ok {
			continue
		}
		for j, r := range def.Replacements {
			ref := ReplacementRef{Object: o.ID, Index: j}
			if !appliedAlready(a, ref) && r.When(g, o.ID, a) {
				out = append(out, ref)
			}
		}
	}
	return out
}

func appliedAlready(a Action, ref ReplacementRef) bool {
	for _, x := range a.Applied {
		if x == ref {
			return true
		}
	}
	return false
}

func (g *Game) replacement(ref ReplacementRef) Replacement {
	def, ok := g.cards.find(g.Object(ref.Object).Card)
	if !ok || ref.Index >= len(def.Replacements) {
		panic(fmt.Sprintf("rulesengine: no replacement %+v", ref))
	}
	return def.Replacements[ref.Index]
}

// affected is the player who orders replacements for a (R-ABIL-33).
func (g *Game) affected(a Action) PlayerID {
	if a.Player == NoPlayer {
		return g.Turn.Active
	}
	return a.Player
}

// processQueue handles the first queued action: it asks for the order of
// several replacements, applies a single one, or performs the action.
func (g *Game) processQueue() {
	a := g.Queue[0]
	refs := g.replacementsFor(a)
	switch {
	case len(refs) > 1:
		opts := make([]string, len(refs))
		for i, r := range refs {
			opts[i] = g.replacement(r).Text
		}
		g.Choices = refs
		g.Waiting = Prompt{Kind: PromptChooseReplacement, Player: g.affected(a), Options: opts}
	case len(refs) == 1:
		g.applyReplacement(refs[0])
	default:
		g.Queue = g.Queue[1:]
		g.perform(a)
	}
}

// applyReplacement rewrites the first queued action with one replacement.
func (g *Game) applyReplacement(ref ReplacementRef) {
	a := g.Queue[0]
	r := g.replacement(ref)
	out := r.Do(g, ref.Object, a)
	for i := range out {
		out[i].Applied = append(append([]ReplacementRef(nil), a.Applied...), ref)
	}
	g.emit(Event{Kind: EvReplaced, Player: g.affected(a), Object: ref.Object, Card: g.Object(ref.Object).Card, Text: r.Text})
	g.Queue = append(out, g.Queue[1:]...)
}

// perform makes an action happen.
func (g *Game) perform(a Action) {
	switch a.Kind {
	case ActGainCents:
		g.Players[a.Player].Cents += a.Amount
		g.emit(Event{Kind: EvGainedCents, Player: a.Player, Amount: a.Amount})
	case ActLoseCents:
		n := min(a.Amount, g.Players[a.Player].Cents) // R-MECH-42
		g.Players[a.Player].Cents -= n
		g.emit(Event{Kind: EvLostCents, Player: a.Player, Amount: n})
	case ActLoot:
		g.loot(a.Player, a.Amount)
	}
}
