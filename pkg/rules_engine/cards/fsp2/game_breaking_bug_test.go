package fsp2

import "testing"

func TestGameBreakingBug(t *testing.T) {
	tb := table(t, seat("isaac", "game_breaking_bug", "mystery_sack"), seat("cain"))
	tb.G.ForceRolls(1)
	tb.Activate(0, "mystery_sack", 0, "mystery_sack")
	if hasItem(tb.G, 0, "mystery_sack") {
		t.Error("on a 1 the sack may be rerolled")
	}
}
