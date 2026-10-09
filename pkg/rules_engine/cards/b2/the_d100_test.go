package b2

import "testing"

func TestTheD100(t *testing.T) {
	tb := itemTable(t, items("the_d100"))
	tb.G.ForceRolls(6)
	tb.Activate(0, "the_d100", 0)
	if a := tb.G.PlayerATK(0); a != 2 {
		t.Errorf("ATK = %d, want 2", a)
	}
}
