package fsp2

import "testing"

func TestHeadOfKrampus(t *testing.T) {
	tb := armed(t, items("head_of_krampus"), "gurdy", "monstro")
	tb.G.ForceRolls(5)
	tb.Activate(0, "head_of_krampus", 0, "ask.rest_in_slot_order")
	if d := tb.G.Object(topID(tb, 0)).Damage; d != 1 {
		t.Errorf("4-6: Gurdy damage = %d, want 1", d)
	}
	tb2 := armed(t, items("head_of_krampus"), "gurdy", "monstro")
	tb2.G.ForceRolls(2)
	tb2.Activate(0, "head_of_krampus", 0)
	if tb2.G.Players[0].Damage != 1 || tb2.G.Players[1].Damage != 1 {
		t.Error("1-3: each player takes 1")
	}
}
