package b2

import "testing"

func TestTheMidasTouch(t *testing.T) {
	tb := itemTable(t, items("the_midas_touch"))
	tb.Attack("fly", 6)
	if c := tb.G.Players[0].Cents; c != 4 { // 3 + the fly's 1¢
		t.Errorf("cents = %d, want 4", c)
	}
}
