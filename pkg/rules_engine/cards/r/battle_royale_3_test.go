package r

import "testing"

func TestBattleRoyale3(t *testing.T) {
	tb := roomTable(t, "battle_royale_3")
	if dc, ok := tb.G.PlayerAttackDC(1); !ok || dc != 5 {
		t.Errorf("Cain's DC = %d, %v; want 5", dc, ok)
	}
}
