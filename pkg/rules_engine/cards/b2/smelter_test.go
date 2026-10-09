package b2

import "testing"

func TestSmelter(t *testing.T) {
	tb := itemTable(t, items("smelter"), "a_dime")
	tb.Activate(0, "smelter", 0, "a_dime")
	if c := tb.G.Players[0].Cents; c != 3 || !inDiscard(tb.G, "a_dime") {
		t.Errorf("cents %d", c)
	}
	if rule := tb.Refused(0, "smelter", 0); rule != "R-ABIL-07" && rule != "R-ABIL-06" {
		t.Errorf("empty hand refused by %s", rule)
	}
}
