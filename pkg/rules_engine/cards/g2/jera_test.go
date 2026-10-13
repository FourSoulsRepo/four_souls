package g2

import "testing"

func TestJera(t *testing.T) {
	tb := lootTable(t, "jera", "bomb", "bomb")
	tb.Play(0, "jera")
	if n := len(tb.G.Players[0].Hand); n != 4 {
		t.Errorf("hand = %d, want 4: 2 cards in hand, loot 2", n)
	}
}
