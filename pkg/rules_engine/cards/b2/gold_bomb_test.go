package b2

import "testing"

func TestGoldBomb(t *testing.T) {
	tb := lootTable(t, "gold_bomb")
	tb.Play(0, "gold_bomb", foe)
	if !tb.G.Players[1].Dead {
		t.Error("3 damage did not kill a 2 HP player")
	}
}
