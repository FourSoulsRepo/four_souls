package r

import "testing"

func TestHolyBrain(t *testing.T) {
	tb := fight(t, "holy_brain")
	tb.Attack("holy_brain", 3, 2) // a 3 hits; each monster has -1 DC: the 2 hits too
	if tb.G.Players[0].Damage != 0 || top(tb, 0) == "holy_brain" {
		t.Errorf("both rolls hit and kill it; damage %d", tb.G.Players[0].Damage)
	}
}
