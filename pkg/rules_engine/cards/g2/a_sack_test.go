package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestASack(t *testing.T) {
	tb := lootTable(t, "a_sack")
	tb.Play(0, "a_sack")
	if n := len(tb.G.Players[0].Hand); n != 3 {
		t.Errorf("hand = %d, want 3", n)
	}
	if !inDiscard(tb.G, engine.LootDeck, "a_sack") {
		t.Error("A Sack is not in the loot discard")
	}
}
