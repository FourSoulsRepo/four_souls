package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestXIITheHangedMan(t *testing.T) {
	tb := lootTable(t, "xii_the_hanged_man")
	treasure := tb.G.DeckTop(engine.TreasureDeck, 1)[0]
	loot := tb.G.DeckTop(engine.LootDeck, 1)[0]
	label := "put " + string(tb.G.Object(treasure).Card) + " on the bottom"
	keep := func(id engine.ObjectID) string { return "keep " + string(tb.G.Object(id).Card) + " on top" }
	monster := tb.G.DeckTop(engine.MonsterDeck, 1)[0]
	tb.Play(0, "xii_the_hanged_man", label, keep(loot), keep(monster))
	if tb.G.Decks[engine.TreasureDeck][0] != treasure {
		t.Error("the treasure card is not on the bottom")
	}
	if n := len(tb.G.Players[0].Hand); n != 2 {
		t.Errorf("hand %d, want 2", n)
	}
}
