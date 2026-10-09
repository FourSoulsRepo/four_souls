package b2

import "testing"

func TestShinyRock(t *testing.T) {
	tb := itemTable(t, items("shiny_rock", "mystery_sack"))
	tb.G.ForceRolls(6)
	tb.Activate(0, "mystery_sack", 0)
	if c := tb.G.Players[0].Cents; c != 1 {
		t.Errorf("cents = %d, want 1", c)
	}
}
