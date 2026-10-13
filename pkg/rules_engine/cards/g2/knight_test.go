package g2

import "testing"

func TestKnight(t *testing.T) {
	tb := fight(t, "knight", "fly")
	tb.G.ForceRolls(1, 1) // two 1s: Isaac takes 1 each, and deals 1 each
	tb.Do(attackIntent(tb), "knight", "isaac")
	if top(tb, 0) == "knight" {
		t.Error("the knight should die from the two 1s")
	}
}
