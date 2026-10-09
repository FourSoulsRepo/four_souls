package b2

import "testing"

func TestMomsDeadHand(t *testing.T) {
	tb := slayTable(t, "moms_dead_hand", seat("isaac"), seat("cain", "breakfast"))
	killWithSixes(t, tb, "breakfast")
	if !hasItem(tb.G, 0, "breakfast") {
		t.Error("nothing stolen")
	}
}
