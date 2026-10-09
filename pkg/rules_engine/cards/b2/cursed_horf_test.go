package b2

import "testing"

func TestCursedHorf(t *testing.T) {
	tb := slayTable(t, "cursed_horf", seat("isaac", "mystery_sack"), seat("cain"))
	tb.G.ForceRolls(2)
	tb.Activate(0, "mystery_sack", 0, "mystery_sack") // dead: the penalty takes the sack
	if !tb.G.Players[0].Dead {
		t.Error("2 damage on a 2 did not kill Isaac")
	}
}
