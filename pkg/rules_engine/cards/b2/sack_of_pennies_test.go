package b2

import "testing"

func TestSackOfPennies(t *testing.T) {
	tb := itemTable(t, items("sack_of_pennies", "mystery_sack"))
	tb.Activate(0, "sack_of_pennies", 0)
	tb.G.ForceRolls(1)
	tb.Activate(0, "mystery_sack", 0, "yes")
	if c := tb.G.Players[0].Cents; c != 1 || !tb.G.Object(tb.Find(0, "sack_of_pennies")).Charged {
		t.Errorf("cents %d, recharged %v", c, tb.G.Object(tb.Find(0, "sack_of_pennies")).Charged)
	}
}
