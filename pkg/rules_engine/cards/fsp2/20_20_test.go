package fsp2

import "testing"

func Test2020(t *testing.T) {
	tb := armed(t, items("20_20"), "leech", "fatty")
	tb.G.ForceRolls(2)
	declareAttack(tb, "leech")
	tb.Activate(0, "20_20", 0, "ask.roll_of 2", "2") // 2 + 2 = 4 hits the leech (DC 4)
	if top(tb, 0) == "leech" {
		t.Error("the leech should die to the changed roll")
	}
}
