package b2

import "testing"

func TestGodhead(t *testing.T) {
	tb := rollTable(t, 3, seat("isaac", "godhead"), seat("cain"))
	tb.Pass(1)
	if got := resolvedRoll(t, tb.Activate(0, "godhead", 0, "ask.roll_of 3", "6")); got != 6 {
		t.Errorf("resolved as %d, want 6", got)
	}
}
