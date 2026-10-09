package b2

import "testing"

func TestCursedMomsHand(t *testing.T) {
	tb := slayTable(t, "cursed_moms_hand", seat("isaac", "mystery_sack"), seat("cain"))
	tb.G.ForceRolls(6)
	tb.Activate(0, "mystery_sack", 0)
	if !tb.G.Turn.EndDeclared {
		t.Error("the turn did not end")
	}
}
