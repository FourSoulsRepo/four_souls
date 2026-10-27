package b2

import "testing"

func TestSpoonBender(t *testing.T) {
	tb := rollTable(t, 3, seat("isaac", "spoon_bender"), seat("cain"))
	tb.Pass(1)
	if got := resolvedRoll(t, tb.Activate(0, "spoon_bender", 0, "ask.roll_of 3")); got != 4 {
		t.Errorf("resolved as %d, want 4", got)
	}
}
