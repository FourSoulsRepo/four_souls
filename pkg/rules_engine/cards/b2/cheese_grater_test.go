package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCheeseGrater(t *testing.T) {
	tb := itemTable(t, items("cheese_grater", "mystery_sack"))
	tb.G.ForceRolls(6)
	tb.Activate(0, "mystery_sack", 0, "ask.deck.treasure", "card.cheese_grater.discard")
	if len(tb.G.Discards[engine.TreasureDeck]) != 1 {
		t.Error("the top treasure card was not discarded")
	}
}
