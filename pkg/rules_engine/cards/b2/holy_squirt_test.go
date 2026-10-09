package b2

import "testing"

func TestHolySquirt(t *testing.T) {
	tb := slayTable(t, "holy_squirt", seat("isaac", "mystery_sack"), seat("cain"))
	tb.G.ForceRolls(5)
	tb.Activate(0, "mystery_sack", 0)
	if h := len(tb.G.Players[0].Hand); h != 1 {
		t.Errorf("hand = %d, want 1", h)
	}
}
