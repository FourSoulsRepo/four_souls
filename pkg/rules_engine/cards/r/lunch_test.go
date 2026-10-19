package r

import "testing"

func TestLunch(t *testing.T) {
	tb := table2(t, "lunch")
	if tb.G.PlayerHP(0) != 3 {
		t.Error("+1 HP")
	}
}
