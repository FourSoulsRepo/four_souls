package r

import "testing"

func TestBrownie(t *testing.T) {
	tb := fight(t, "brownie")
	tb.Attack("brownie", 5, 4, 3, 2, 1, 1) // DC 5, then 4, 3, 2, 1, 1
	if tb.G.SoulValue(0) != 1 {
		t.Error("each hit lowers its DC: six rolls kill it")
	}
}
