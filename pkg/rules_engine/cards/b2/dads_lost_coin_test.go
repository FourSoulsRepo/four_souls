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
	tb.Activate(0, "mystery_sack", 0, "No")
	if h := len(tb.G.Players[0].Hand); h != 1 {
		t.Errorf("hand = %d, want 1 (the 1 stands)", h)
	}
}

// The reroll is a 1 again: the coin triggers again (19.66).
func TestDadsLostCoinTwice(t *testing.T) {
	tb := itemTable(t, items("dads_lost_coin", "mystery_sack"))
	tb.G.ForceRolls(1, 1, 3)
	tb.Activate(0, "mystery_sack", 0, "yes", "yes")
	if c := tb.G.Players[0].Cents; c != 4 {
		t.Errorf("cents = %d, want 4 (rerolled twice to 3)", c)
	}
}
