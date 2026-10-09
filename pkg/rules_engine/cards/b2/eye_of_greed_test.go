package b2

import "testing"

func TestEyeOfGreed(t *testing.T) {
	tb := itemTable(t, items("eye_of_greed", "mystery_sack"))
	tb.G.ForceRolls(5)
	tb.Activate(0, "mystery_sack", 0)
	if c := tb.G.Players[0].Cents; c != 3 {
		t.Errorf("cents = %d, want 3", c)
	}
}
