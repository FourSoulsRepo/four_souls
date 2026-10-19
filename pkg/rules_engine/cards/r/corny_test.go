package r

import "testing"

func TestCorny(t *testing.T) {
	tb := fight(t, "corny")
	tb.Attack("corny", 6)
	if tb.G.Players[0].Damage != 0 || tb.G.Players[1].Damage+tb.G.Players[2].Damage != 2 {
		t.Error("2 damage to a non-active player at random")
	}
}
