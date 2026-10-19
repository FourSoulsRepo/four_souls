package r

import "testing"

func TestBattleRoyale2(t *testing.T) {
	tb := roomTable(t, "battle_royale_2")
	if dc, ok := tb.G.PlayerAttackDC(1); !ok || dc != 4 {
		t.Errorf("Cain's DC = %d, %v; want 4", dc, ok)
	}
}
