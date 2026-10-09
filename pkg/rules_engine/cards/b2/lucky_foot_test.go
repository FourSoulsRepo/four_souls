package b2

import "testing"

func TestLuckyFoot(t *testing.T) {
	tb := rollTable(t, 2, seat("isaac", "lucky_foot"), seat("cain"))
	tb.Pass(1)
	if got := resolvedRoll(t, tb.Activate(0, "lucky_foot", 0, "roll of 2", "2")); got != 4 {
		t.Errorf("resolved as %d, want 4", got)
	}
}
