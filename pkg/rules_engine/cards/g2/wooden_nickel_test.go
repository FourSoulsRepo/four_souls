package g2

import "testing"

func TestWoodenNickel(t *testing.T) {
	tb := table(t, seat("the_keeper", "wooden_nickel"), seat("cain"))
	tb.G.ForceRolls(4)
	tb.Activate(0, "wooden_nickel", 0, foe)
	if got := tb.G.Players[1].Cents; got != 4 {
		t.Errorf("Cain's cents = %d, want 4", got)
	}
	if got := tb.G.Players[0].Cents; got != 0 {
		t.Errorf("Keeper's cents = %d, want 0", got)
	}
}
