package fsp2

import "testing"

func TestPHD(t *testing.T) {
	tb := table(t, seat("isaac", "phd", "mystery_sack"), seat("cain"))
	tb.G.ForceRolls(2)
	tb.Activate(0, "mystery_sack", 0, "+1") // a 2 becomes a 3, once (R-DICE-10)
	if c := tb.G.Players[0].Cents; c != 4 {
		t.Errorf("cents = %d: the sack's 3 gives 4¢", c)
	}
}
