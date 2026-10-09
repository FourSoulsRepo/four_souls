package b2

import "testing"

func TestBabyHaunt(t *testing.T) {
	tb := itemTable(t, items("baby_haunt"))
	fly, _ := tb.G.Monsters[0].TopOf()
	if dc := tb.G.Evasion(fly); dc != 3 {
		t.Errorf("fly DC on your turn = %d, want 3", dc)
	}
}
