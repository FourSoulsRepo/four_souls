package b2

import "testing"

func TestMomsBox(t *testing.T) {
	tb := itemTable(t, items("moms_box", "mystery_sack"), "a_dime")
	tb.G.ForceRolls(4)
	tb.Activate(0, "mystery_sack", 0, "yes", "a_dime")
	if h := handCards(tb.G, 0); len(h) != 1 || h[0] == "a_dime" {
		t.Errorf("hand %v", h)
	}
}
