package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// Requiem cards are played with the Base Game; its reprints need it.
var sets = []engine.CardSet{b2.Set, Set}

// seat is a player with a character and items.
var seat = enginetest.Seat

// setup builds a table from a full setup.
func setup(t *testing.T, s engine.SituationSetup) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, s, sets...)
}

// lootTable: player 0 (Isaac) holds the cards, player 1 is Cain; a fly
// (DC 2) and a leech (DC 4) are in the monster slots.
func lootTable(t *testing.T, hand ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: hand}, seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
}

// inDiscard reports whether the deck's discard holds the card.
func inDiscard(g *engine.Game, d engine.DeckKind, card engine.CardRef) bool {
	for _, id := range g.Discards[d] {
		if g.Object(id).Card == card {
			return true
		}
	}
	return false
}

const (
	me  = "player 1 (isaac)"
	foe = "player 2 (cain)"
)

// hasItem reports whether player p controls an item with this card.
func hasItem(g *engine.Game, p int, card engine.CardRef) bool {
	for _, id := range g.Players[p].InPlay {
		if o := g.Object(id); o.Card == card && o.Role == engine.RoleItem {
			return true
		}
	}
	return false
}

// handHas reports whether player p holds a loot card with this card.
func handHas(g *engine.Game, p int, card engine.CardRef) bool {
	for _, id := range g.Players[p].Hand {
		if g.Object(id).Card == card {
			return true
		}
	}
	return false
}

// top returns the card on top of monster slot i.
func top(tb *enginetest.Table, i int) engine.CardRef {
	id, _ := tb.G.Monsters[i].TopOf()
	return tb.G.Object(id).Card
}

// settleAny settles, answering every question with its first option.
func settleAny(tb *enginetest.Table) {
	tb.T.Helper()
	for range 200 {
		tb.ThrowPending()
		w := tb.G.Prompt()
		switch {
		case w.Kind == engine.PromptChoose:
			tb.Choose(w.Options[0])
			continue
		case w.Kind == engine.PromptDiscard:
			if _, err := tb.G.Apply(engine.Intent{Player: w.Player, Kind: engine.IntentDiscard, Objects: tb.G.Players[w.Player].Hand[:w.Count]}); err != nil {
				tb.T.Fatal(err)
			}
			continue
		case len(tb.G.Stack) == 0 && w.Kind == engine.PromptPriority && !tb.G.Attack.On:
			return
		case w.Kind == engine.PromptPriority:
			tb.Pass(int(w.Player))
			continue
		}
		return
	}
	tb.T.Fatal("settleAny: the stack never emptied")
}
