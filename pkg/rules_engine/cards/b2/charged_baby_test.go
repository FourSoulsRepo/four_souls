package b2

import "testing"

func TestChargedBaby(t *testing.T) {
	tb := itemTable(t, items("charged_baby", "mystery_sack"))
	tb.G.ForceRolls(2)
	tb.Activate(0, "mystery_sack", 0, "mystery_sack")
	if !tb.G.Object(tb.Find(0, "mystery_sack")).Charged {
		t.Error("the sack was not recharged")
	}
}
