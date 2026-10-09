package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestXIXTheSun(t *testing.T) {
	tb := lootTable(t, "xix_the_sun")
	tb.Play(0, "xix_the_sun")
	if c := tb.G.Object(tb.G.Decks[engine.LootDeck][0]).Card; c != "xix_the_sun" {
		t.Errorf("bottom of the loot deck is %s", c)
	}
	tb.EndTurn()
	if tb.G.Turn.Active != 0 {
		t.Error("no extra turn")
	}
}
