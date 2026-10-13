package fsp2

import "testing"

func TestCursedTumor(t *testing.T) {
	tb := armed(t, items("mystery_sack"), "cursed_tumor", "gurdy")
	tb.G.ForceRolls(4)
	tb.Activate(0, "mystery_sack", 0)
	if tb.G.Players[0].Damage != 1 || tb.G.Players[1].Damage != 1 {
		t.Error("a 4: each player takes 1")
	}
}
