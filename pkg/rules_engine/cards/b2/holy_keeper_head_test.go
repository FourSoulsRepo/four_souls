package b2

import "testing"

func TestHolyKeeperHead(t *testing.T) {
	tb := slayTable(t, "holy_keeper_head", seat("isaac", "mystery_sack"), seat("cain"))
	tb.G.ForceRolls(4)
	tb.Activate(0, "mystery_sack", 0)
	if c := tb.G.Players[0].Cents; c != 6 {
		t.Errorf("cents = %d, want 6 (4 + 2)", c)
	}
}
