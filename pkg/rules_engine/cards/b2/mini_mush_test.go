package b2

import "testing"

func TestMiniMush(t *testing.T) {
	tb := rollTable(t, 5, seat("isaac", "mini_mush"), seat("cain"))
	tb.Pass(1)
	if got := resolvedRoll(t, tb.Activate(0, "mini_mush", 0, "roll of 5", "-2")); got != 3 {
		t.Errorf("resolved as %d, want 3", got)
	}
}
