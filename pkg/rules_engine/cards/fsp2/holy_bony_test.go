package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHolyBony(t *testing.T) {
	tb := armed(t, items("mystery_sack"), "holy_bony", "gurdy")
	tb.G.PutOnDeck(engine.LootDeck, "bomb", "two_cents", "four_cents") // four_cents on top
	tb.G.ForceRolls(6)
	tb.Activate(0, "mystery_sack", 0, "loot deck", "bomb", "two_cents", "four_cents")
	if top := tb.G.DeckTop(engine.LootDeck, 1); tb.G.Object(top[0]).Card != "bomb" {
		t.Error("the roller put the bomb on top")
	}
}
