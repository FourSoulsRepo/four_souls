package b2

import "testing"

func TestHolyDip(t *testing.T) {
	tb := slayTable(t, "holy_dip", seat("isaac", "mystery_sack"), seat("cain"))
	tb.G.ForceRolls(1)
	tb.Activate(0, "mystery_sack", 0)
	if c := tb.G.Players[0].Cents; c != 1 {
		t.Errorf("cents = %d, want 1", c)
	}
}
