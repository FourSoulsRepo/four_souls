package rulesengine

import (
	"fmt"
	"strconv"
)

// ActionKind is a change to the game that cards may rewrite (A-05).
type ActionKind int

// The action kinds; more come with combat and effect blocks.
const (
	ActGainCents          ActionKind = iota // R-MECH-36
	ActLoseCents                            // R-MECH-42
	ActLoot                                 // R-CARD-06
	ActGainTreasure                         // R-CARD-02
	ActBecomeSoul                           // a dead monster becomes a soul (R-DEATH-08)
	ActDiscardObject                        // put an object into its discard (R-ZONE-06)
	ActRefillSlots                          // R-SHOP-06
	ActPenaltyItem                          // death penalty: destroy an item (R-DEATH-14)
	ActPenaltyLoot                          // death penalty: discard a loot card
	ActDeactivateTaps                       // death penalty: deactivate ↷ objects
	ActAsk                                  // an Ask effect asks its questions (R-ABIL-05)
	ActStealCents                           // Player steals Amount¢ From a player (R-MECH-38)
	ActChooseStartingItem                   // Eden-style start-of-game choice (R-SETUP-09)
	ActPenaltyDone                          // the death penalty is paid
	ActAddCounters                          // put Amount counters on Object (Bum-bo levels)
)

// Action is a pending change. It sits in the queue, may be rewritten by
// replacement effects, and only then happens (R-ABIL-29).
type Action struct {
	Kind   ActionKind `json:"kind"`
	Player PlayerID   `json:"player"`
	Object ObjectID   `json:"object,omitempty"`
	Amount int        `json:"amount"`
	From   PlayerID   `json:"from,omitempty"` // the other player, e.g. of a steal
	Ask    *Asking    `json:"ask,omitempty"`
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
	for _, id := range g.inPlay() {
		o := g.Object(id)
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
		g.ask(Choice{Purpose: ChooseReplacement, Player: g.affected(a), Rule: "R-ABIL-33", Replacements: refs}, opts)
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
		n := a.Amount << g.bonus(StatLootDouble, a.Player, 0) // Two of Clubs
		if g.Compost {
			g.Compost = false
			g.lootFromDiscard(a.Player, n)
			return
		}
		g.loot(a.Player, n)
	case ActAddCounters:
		if o := g.Object(a.Object); o.Zone.Kind == ZoneInPlay {
			o.addCounters("", a.Amount)
			g.emit(Event{Kind: EvCounters, Player: o.Controller, Object: a.Object, Card: o.Card, Amount: a.Amount})
		}
	case ActGainTreasure:
		g.gainTreasure(a.Player, a.Amount)
	case ActBecomeSoul:
		nid := g.move(a.Object, Zone{Kind: ZoneInPlay}, a.Player)
		g.Object(nid).Role = RoleSoul
		g.Players[a.Player].InPlay = append(g.Players[a.Player].InPlay, nid)
		g.emit(Event{Kind: EvGainedSoul, Player: a.Player, Object: nid, Card: g.Object(nid).Card})
	case ActDiscardObject:
		if deck, ok := g.kindOf(a.Object).Deck(); ok {
			g.discard(a.Object, deck)
		}
	case ActRefillSlots:
		g.refillSlots()
	case ActPenaltyItem:
		var items []ObjectID
		for _, id := range g.Players[a.Player].InPlay {
			if g.Object(id).Role == RoleItem && !g.Eternal(id) {
				items = append(items, id)
			}
		}
		if len(items) > 0 {
			chooser := a.Player
			if s := g.shadowOf(a.Player); s != NoPlayer {
				chooser = s
			}
			g.ask(Choice{Purpose: ChoosePenaltyItem, Player: chooser, Owner: a.Player, To: NoPlayer, Rule: "R-DEATH-14", Objects: items}, g.labels(items))
		}
	case ActPenaltyLoot:
		if hand := g.Players[a.Player].Hand; len(hand) > 0 {
			g.ask(Choice{Purpose: ChoosePenaltyLoot, Player: a.Player, Owner: a.Player, To: g.shadowOf(a.Player), Rule: "R-DEATH-14", Objects: append([]ObjectID(nil), hand...)}, g.labels(hand))
		}
	case ActDeactivateTaps:
		pl := g.Players[a.Player]
		g.Object(pl.Character).Charged = false
		for _, id := range pl.InPlay {
			if g.Object(id).Role == RoleItem && g.def(id).Tap {
				g.Object(id).Charged = false
			}
		}
		g.emit(Event{Kind: EvDeactivated, Player: a.Player})
	case ActAsk:
		g.continueAsk(a.Ask)
	case ActStealCents:
		n := min(a.Amount, g.Players[a.From].Cents)
		g.Players[a.From].Cents -= n
		g.Players[a.Player].Cents += n
		g.emit(Event{Kind: EvStole, Player: a.Player, Amount: n, Text: strconv.Itoa(int(a.From))})
	case ActChooseStartingItem:
		g.askStartingItem(a.Player, a.Amount)
	case ActPenaltyDone:
		g.emit(Event{Kind: EvPenaltyPaid, Player: a.Player})
	}
}

// gainTreasure puts the top n treasure cards into play under p (R-CARD-02).
func (g *Game) gainTreasure(p PlayerID, n int) {
	for range n {
		id, ok := g.drawTop(TreasureDeck)
		if !ok {
			return
		}
		nid := g.move(id, Zone{Kind: ZoneInPlay}, p)
		o := g.Object(nid)
		g.enterAsItem(p, nid)
		g.emit(Event{Kind: EvGainedTreasure, Player: p, Object: nid, Card: o.Card})
	}
}
