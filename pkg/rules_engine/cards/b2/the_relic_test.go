package b2

import "testing"

func TestTheRelic(t *testing.T) {
	tb := itemTable(t, items("the_relic", "mystery_sack"))
	tb.G.ForceRolls(1)
	tb.Activate(0, "mystery_sack", 0)
	if h := len(tb.G.Players[0].Hand); h != 2 { // the relic and the sack
		t.Errorf("hand = %d, want 2", h)
	}
}
