package b2

import "testing"

func TestBox(t *testing.T) {
	tb := itemTable(t, items("box"), "a_penny_6", "a_penny_6", "a_penny_6")
	tb.Activate(0, "box", 0)
	for range 3 {
		tb.Play(0, "a_penny_6")
	}
	if c := tb.G.Players[0].Cents; c != 3 || hasItem(tb.G, 0, "box") {
		t.Errorf("cents %d, box still there %v", c, hasItem(tb.G, 0, "box"))
	}
}
