package b2

import "testing"

func TestDadsLostCoin(t *testing.T) {
	tb := itemTable(t, items("dads_lost_coin", "mystery_sack"))
	tb.G.ForceRolls(1, 3) // the 1 would resolve: forced reroll gives 3
	tb.Activate(0, "mystery_sack", 0, "yes")
	if c := tb.G.Players[0].Cents; c != 4 {
		t.Errorf("cents = %d, want 4 (rerolled to 3)", c)
	}
}

func TestDadsLostCoinDeclined(t *testing.T) {
	tb := itemTable(t, items("dads_lost_coin", "mystery_sack"))
	tb.G.ForceRolls(1)
	tb.Activate(0, "mystery_sack", 0, "no")
	if h := len(tb.G.Players[0].Hand); h != 1 {
		t.Errorf("hand = %d, want 1 (the 1 stands)", h)
	}
}
