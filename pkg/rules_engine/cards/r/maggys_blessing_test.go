package r

import "testing"

func TestMaggysBlessing(t *testing.T) {
	tb := roomTable(t, "maggys_blessing")
	if got := tb.G.PlayerHP(1); got != 3 {
		t.Errorf("Cain's HP = %d, want 3", got)
	}
}
