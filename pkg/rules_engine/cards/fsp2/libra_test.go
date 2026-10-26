package fsp2

import "testing"

func TestLibra(t *testing.T) {
	tb := armed(t, items("libra"), "fly", "fatty")
	tb.G.ForceRolls(1)
	declareAttack(tb, "fly")
	tb.Activate(0, "libra", 0, "ask.roll_of 1")
	if top(tb, 0) == "fly" {
		t.Error("the 1 became a 3 and hits the fly (DC 2)")
	}
}
