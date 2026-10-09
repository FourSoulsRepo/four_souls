package b2

import "testing"

func TestJawbone(t *testing.T) {
	tb := itemTable(t, items("jawbone"))
	tb.G.Players[1].Cents = 2
	tb.Activate(0, "jawbone", 0, foe)
	if a, b := tb.G.Players[0].Cents, tb.G.Players[1].Cents; a != 2 || b != 0 {
		t.Errorf("cents %d and %d, want 2 and 0 (all they had)", a, b)
	}
}
