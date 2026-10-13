package fsp2

import "testing"

func TestSucker(t *testing.T) {
	tb := fight(t, "sucker", "fly")
	tb.G.ForceRolls(6)
	tb.Do(attackIntent(tb), "sucker")
	if tb.G.Players[0].Damage != 1 || tb.G.Players[1].Damage != 1 {
		t.Error("each player takes 1 when it dies")
	}
}
