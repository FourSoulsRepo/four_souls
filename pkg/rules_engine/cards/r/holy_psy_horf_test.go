package r

import "testing"

func TestHolyPsyHorf(t *testing.T) {
	tb := fight(t, "holy_psy_horf")
	tb.EndTurn()
	if len(tb.G.Players[0].Hand) != 1 {
		t.Error("at the end of the turn the active player loots 1")
	}
}
