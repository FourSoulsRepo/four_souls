package b2

import "testing"

func TestGuppysCollar(t *testing.T) {
	tb := itemTable(t, items("guppys_collar"), "xiii_death")
	tb.G.ForceRolls(2)
	tb.Play(0, "xiii_death", me)
	if tb.G.Players[0].Dead {
		t.Error("a 2 should prevent the death")
	}
}
