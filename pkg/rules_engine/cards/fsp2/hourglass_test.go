package fsp2

import "testing"

func TestHourglass(t *testing.T) {
	tb := table(t, seat("isaac", "hourglass", "mystery_sack"), seat("cain", "boomerang"))
	tb.G.ForceRolls(2)
	tb.Activate(0, "mystery_sack", 0, "boomerang")
	if tb.G.Object(tb.Find(1, "boomerang")).Charged {
		t.Error("on a 2 an item may be deactivated")
	}
}
