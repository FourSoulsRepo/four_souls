package b2

import "testing"

func TestThePolaroid(t *testing.T) {
	tb := itemTable(t, items("the_polaroid"))
	tb.EndTurn()
	if h := len(tb.G.Players[0].Hand); h != 2 {
		t.Errorf("hand = %d, want 2", h)
	}
}
