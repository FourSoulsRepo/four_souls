package r

import "testing"

func TestSamsonsBlessing(t *testing.T) {
	tb := roomTable(t, "samsons_blessing")
	if got := tb.G.PlayerATK(1); got != 2 {
		t.Errorf("Cain's ATK = %d, want 2", got)
	}
}
