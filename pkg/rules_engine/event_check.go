package rulesengine

import (
	"fmt"
	"slices"
	"strings"
)

// Snapshot is what CheckEvents compares: each object's zone, damage
// and charge, and the order of every deck (22.2, TS-09).
type Snapshot struct {
	objects []objectState
	decks   [deckCount][]ObjectID
}

type objectState struct {
	Zone    Zone
	Damage  int
	Charged bool
}

// moveLeg is one move: the card of from became the new object to.
// at is the number of the step's events before it.
type moveLeg struct {
	from, to         ObjectID
	fromZone, toZone Zone
	at               int
}

// changeLog is what Game.move and ShuffleDeck log for CheckEvents.
type changeLog struct {
	moves    []moveLeg
	shuffled [deckCount]bool
}

// Change is one state change in a step that CheckEvents found.
type Change struct {
	Object ObjectID
	Card   CardRef
	// What changed: "move" (a new object; Prev is the old one), "new"
	// (a card made in a zone), "zone" (the same object), "damage",
	// "charge", "deck order", or "shuffle" (a whole deck, Object -1).
	What string
	Prev ObjectID
	// From and To are the zones before and after.
	From, To Zone
	// Before and After are the old and new values of damage or charge.
	Before, After int
	// at: for a move, the number of the step's events before it.
	at int
}

func (c Change) String() string {
	s := fmt.Sprintf("object %d (%s) %s: %s -> %s", c.Object, c.Card, c.What, zoneText(c.From), zoneText(c.To))
	switch c.What {
	case "move":
		s = fmt.Sprintf("object %d (%s) moved from %d: %s -> %s", c.Object, c.Card, c.Prev, zoneText(c.From), zoneText(c.To))
	case "shuffle":
		s = "the " + zoneText(c.To) + " shuffled"
	case "damage", "charge":
		s = fmt.Sprintf("object %d (%s) %s in %s: %d -> %d", c.Object, c.Card, c.What, zoneText(c.To), c.Before, c.After)
	}
	return s
}

// SilentChange is a kind of change that needs no event of its own.
type SilentChange struct {
	Name string
	// Reason says why it is silent. Bug names the 19.N item of a gap
	// that 22.3 fixes; it is empty for a change silent on purpose.
	Reason, Bug string
	Match       func(g *Game, c Change, events []Event) bool
}

// TakeSnapshot records the state CheckEvents compares, and logs moves
// from now on.
func (g *Game) TakeSnapshot() Snapshot {
	s := Snapshot{objects: make([]objectState, len(g.Objects))}
	for i, o := range g.Objects {
		s.objects[i] = objectState{Zone: o.Zone, Damage: o.Damage, Charged: o.Charged}
	}
	for d := range deckCount {
		s.decks[d] = slices.Clone(g.Decks[d])
	}
	g.changes = &changeLog{}
	return s
}

// CheckEvents reports the changes since before that no event names,
// or nil; events are the step's events. Fuzz runs call it after every
// step (22.2).
//
// An event names an object by its Object, Prev, Source or stack item
// source. A move makes a new object: an event must name the new one,
// or the old one after the move. A few events name many objects
// (wideEvents); the changes in silentChanges need no event.
func (g *Game) CheckEvents(before Snapshot, events []Event) error {
	changes := g.Changes(before)
	g.changes = nil
	if len(changes) == 0 {
		return nil
	}
	namedAt := map[ObjectID]int{} // the last event that names the object
	for i, e := range events {
		ids := []ObjectID{e.Object, e.Prev, e.Source}
		if e.Item != nil {
			ids = append(ids, e.Item.Source)
		}
		for _, id := range ids {
			namedAt[id] = i + 1
		}
	}
	var missing []string
	for _, c := range changes {
		if _, ok := namedAt[c.Object]; ok {
			continue
		}
		if c.What == "move" && namedAt[c.Prev] > c.at {
			continue
		}
		if wideNamed(g, c, events) || silent(g, c, events) {
			continue
		}
		missing = append(missing, c.String())
	}
	if len(missing) == 0 {
		return nil
	}
	kinds := make([]string, len(events))
	for i, e := range events {
		kinds[i] = string(e.Kind)
	}
	return fmt.Errorf("no event names %s; events: %s (22.2)", strings.Join(missing, "; "), strings.Join(kinds, " "))
}

// Changes lists the moves since TakeSnapshot, and each other zone,
// damage, charge and deck order change since before.
func (g *Game) Changes(before Snapshot) []Change {
	var out []Change
	log := g.changes
	if log == nil {
		log = &changeLog{}
	}
	made := map[ObjectID]bool{}
	for _, m := range log.moves {
		made[m.to] = true
		out = append(out, Change{Object: m.to, Card: g.Objects[m.to].Card, What: "move", Prev: m.from, From: m.fromZone, To: m.toZone, at: m.at})
	}
	for i := range g.Objects {
		o := &g.Objects[i]
		id := ObjectID(i)
		if i >= len(before.objects) {
			if !made[id] && o.Zone.Kind != ZoneGone {
				out = append(out, Change{Object: id, Card: o.Card, What: "new", From: Zone{Kind: ZoneGone}, To: o.Zone})
			}
			continue
		}
		b := before.objects[i]
		if o.Zone.Kind == ZoneGone {
			continue // a move, listed above
		}
		c := Change{Object: id, Card: o.Card, From: b.Zone, To: o.Zone}
		if b.Zone != o.Zone {
			c.What = "zone"
			out = append(out, c)
		}
		if b.Damage != o.Damage {
			c.What, c.Before, c.After = "damage", b.Damage, o.Damage
			out = append(out, c)
		}
		if b.Charged != o.Charged {
			c.What, c.Before, c.After = "charge", boolInt(b.Charged), boolInt(o.Charged)
			out = append(out, c)
		}
	}
	for d := range deckCount {
		moved := reordered(before.decks[d], g.Decks[d])
		z := DeckZone(d)
		if log.shuffled[d] {
			if len(moved) > 0 { // one change for the whole deck
				out = append(out, Change{Object: -1, What: "shuffle", From: z, To: z})
			}
			continue
		}
		for _, id := range moved {
			out = append(out, Change{Object: id, Card: g.Objects[id].Card, What: "deck order", From: z, To: z})
		}
	}
	return out
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// reordered returns the cards in both before and after whose order
// changed: those off the longest run kept in the same order. A card
// put on the bottom is one; a shuffle moves most.
func reordered(before, after []ObjectID) []ObjectID {
	pos := make(map[ObjectID]int, len(before))
	for i, id := range before {
		pos[id] = i
	}
	var ids []ObjectID
	var seq []int
	for _, id := range after {
		if p, ok := pos[id]; ok {
			ids = append(ids, id)
			seq = append(seq, p)
		}
	}
	keep := longestRising(seq)
	var out []ObjectID
	for i, id := range ids {
		if !keep[i] {
			out = append(out, id)
		}
	}
	return out
}

// longestRising marks one longest rising run of seq (patience sort).
func longestRising(seq []int) []bool {
	var tails []int // index in seq of the smallest tail of each length
	prev := make([]int, len(seq))
	for i, v := range seq {
		n, _ := slices.BinarySearchFunc(tails, v, func(t, v int) int { return seq[t] - v })
		prev[i] = -1
		if n > 0 {
			prev[i] = tails[n-1]
		}
		if n == len(tails) {
			tails = append(tails, i)
		} else {
			tails[n] = i
		}
	}
	keep := make([]bool, len(seq))
	if len(tails) > 0 {
		for i := tails[len(tails)-1]; i >= 0; i = prev[i] {
			keep[i] = true
		}
	}
	return keep
}

func zoneText(z Zone) string {
	names := [...]string{"gone", "deck", "discard", "hand", "in play", "covered", "stack", "outside"}
	s := names[z.Kind]
	switch z.Kind {
	case ZoneDeck, ZoneDiscard:
		s = z.Deck.String() + " " + s
	case ZoneInPlay, ZoneCovered:
		if z.Index != 0 || z.Slot != ShopSlot {
			s += fmt.Sprintf(" slot %d.%d", z.Slot, z.Index)
		}
	default: // the kind alone says where
	}
	return s
}

// WideEvent is an event kind that names many objects at once: all a
// player controls, or all the cards that went into a deck.
type WideEvent struct {
	Kind   EventKind
	Reason string
	Match  func(g *Game, c Change, e Event) bool
}

// wideEvents name many objects, not each one.
var wideEvents = []WideEvent{
	{
		EvRecharged, "the recharge step and a recharge choice charge the player's things (R-TURN-02)",
		func(g *Game, c Change, e Event) bool { return c.What == "charge" && ownEvent(g, c, e) },
	},
	{
		EvDeactivated, "the death penalty deactivates all the player's ↷ things at once (R-DEATH-14)",
		func(g *Game, c Change, e Event) bool { return c.What == "charge" && ownEvent(g, c, e) },
	},
	{
		EvDeckBottom, "the cards that went under a deck are hidden; the event counts them (19.147)",
		func(_ *Game, c Change, e Event) bool {
			return c.To.Kind == ZoneDeck && e.Text == c.To.Deck.String()+" deck"
		},
	},
	{
		EvMovedToDeck, "an empty deck takes back its shuffled discard; the event counts the cards (R-ZONE-05, 19.174)",
		func(_ *Game, c Change, e Event) bool {
			return e.Object == 0 && c.What == "move" && c.From.Kind == ZoneDiscard &&
				c.To == DeckZone(c.From.Deck) && e.Text == c.To.Deck.String()+" deck"
		},
	},
	{
		EvTurnEnded, "every object heals at the end of the turn (R-TURN-13)",
		func(_ *Game, c Change, _ Event) bool { return c.What == "damage" && c.After == 0 },
	},
}

// ownEvent: the event names the player who controls the changed
// object, and no object.
func ownEvent(g *Game, c Change, e Event) bool {
	return e.Object == 0 && e.Player == g.Objects[c.Object].Controller
}

func wideNamed(g *Game, c Change, events []Event) bool {
	for _, w := range wideEvents {
		for _, e := range events {
			if e.Kind == w.Kind && w.Match(g, c, e) {
				return true
			}
		}
	}
	return false
}

// silentChanges need no event: on purpose, or a gap logged as a 19.N
// item that 22.3 fixes (Bug set).
var silentChanges = []SilentChange{
	{
		Name:   "covered",
		Reason: "a card that enters a slot on top covers the one there; its entered_play says so (R-ZONE-10)",
		Match: func(g *Game, c Change, events []Event) bool {
			if c.What != "zone" || c.From.Kind != ZoneInPlay || c.To.Kind != ZoneCovered ||
				c.From.Slot != c.To.Slot || c.From.Index != c.To.Index {
				return false
			}
			for _, e := range events {
				if e.Kind == EvEnteredPlay && int(e.Object) < len(g.Objects) {
					if z := g.Objects[e.Object].Zone; z.Kind == ZoneInPlay && z.Slot == c.To.Slot && z.Index == c.To.Index {
						return true
					}
				}
			}
			return false
		},
	},
	{
		Name:   "minion slot renumbered",
		Reason: "a removed minion slot shifts the later slots' numbers down; no card moves (R-CHAL-11, ADR 028 7.7)",
		Match: func(_ *Game, c Change, _ []Event) bool {
			return c.What == "zone" && c.From.Kind == c.To.Kind && c.From.Slot == MinionSlot &&
				c.To.Slot == MinionSlot && c.To.Index < c.From.Index
		},
	},
	{
		Name:   "deactivated by a card",
		Reason: "Spider Webs, Curse of Impulse and Hourglass deactivate a card with no event",
		Bug:    "19.176",
		Match: func(_ *Game, c Change, _ []Event) bool {
			return c.What == "charge" && c.After == 0
		},
	},
	{
		Name:   "shuffle",
		Reason: "a shuffled deck shows no card and no count change; nothing to draw",
		Match:  func(_ *Game, c Change, _ []Event) bool { return c.What == "shuffle" },
	},
	{
		Name:   "max HP fell",
		Reason: "a lower max HP removes as much damage; not healing, and the HP left stays the same (R-MECH-21)",
		Match: func(_ *Game, c Change, _ []Event) bool {
			return c.What == "damage" && c.After < c.Before && c.To.Kind == ZoneInPlay
		},
	},
}

func silent(g *Game, c Change, events []Event) bool {
	for _, s := range silentChanges {
		if s.Match(g, c, events) {
			return true
		}
	}
	return false
}

// SilentChanges lists the changes that need no event, for tests.
func SilentChanges() []SilentChange { return slices.Clone(silentChanges) }
