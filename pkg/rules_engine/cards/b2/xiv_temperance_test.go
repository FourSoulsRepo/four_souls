package b2

import "testing"

func TestXIVTemperance(t *testing.T) {
	tb := lootTable(t, "xiv_temperance")
	tb.Play(0, "xiv_temperance", "Take 1 damage and gain 4¢.")
	if d, c := tb.G.Players[0].Damage, tb.G.Players[0].Cents; d != 1 || c != 4 {
		t.Errorf("damage %d cents %d, want 1 and 4", d, c)
	}
	tb = lootTable(t, "xiv_temperance")
	tb.Play(0, "xiv_temperance", "Take 2 damage and gain 8¢.")
	// 2 damage kills Isaac: 8¢ minus the 1¢ death penalty (R-DEATH-14).
	if c := tb.G.Players[0].Cents; c != 7 || !tb.G.Players[0].Dead {
		t.Errorf("cents %d dead %v, want 7 and dead", c, tb.G.Players[0].Dead)
	}
}
