package b2

import "testing"

func TestMysterySack(t *testing.T) {
	tb := itemTable(t, items("mystery_sack"))
	tb.G.ForceRolls(3)
	tb.Activate(0, "mystery_sack", 0)
	if c := tb.G.Players[0].Cents; c != 4 {
		t.Errorf("cents = %d, want 4", c)
	}
}
