package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

const (
	lookHand = "Look at a player's hand. You may swap a card from your hand with one of theirs."
	lootPut  = "Loot 1, then put a card from your hand on top of the loot deck."
)

func handTable(t *testing.T) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "lilith", Items: []engine.CardRef{"incubus"}, Hand: []engine.CardRef{"a_nickel"}},
		{Character: "cain", Hand: []engine.CardRef{"a_dime"}},
	}}, Set)
}

func handCards(g *engine.Game, p int) []engine.CardRef {
	var out []engine.CardRef
	for _, id := range g.Players[p].Hand {
		out = append(out, g.Object(id).Card)
	}
	return out
}

func TestIncubusSwaps(t *testing.T) {
	tb := handTable(t)
	tb.Activate(0, "incubus", 0, lookHand, "player 2 (cain)", "a_dime", "a_nickel")
	if a, b := handCards(tb.G, 0), handCards(tb.G, 1); len(a) != 1 || a[0] != "a_dime" || len(b) != 1 || b[0] != "a_nickel" {
		t.Errorf("hands %v and %v, want swapped", a, b)
	}
}

func TestIncubusMayNotSwap(t *testing.T) {
	tb := handTable(t)
	ev := tb.Activate(0, "incubus", 0, lookHand, "player 2 (cain)", "don't swap")
	if a := handCards(tb.G, 0); len(a) != 1 || a[0] != "a_nickel" {
		t.Errorf("hand %v changed", a)
	}
	if !enginetest.Has(ev, engine.EvLookedAt, 0) {
		t.Error("no look at the hand")
	}
}

func TestIncubusLootsAndPutsBack(t *testing.T) {
	tb := handTable(t)
	tb.Activate(0, "incubus", 0, lootPut, "a_nickel")
	if n := len(tb.G.Players[0].Hand); n != 1 {
		t.Errorf("hand has %d cards, want 1", n)
	}
	if top := tb.G.DeckTop(engine.LootDeck, 1); tb.G.Object(top[0]).Card != "a_nickel" {
		t.Error("A Nickel! is not on top of the loot deck")
	}
}
