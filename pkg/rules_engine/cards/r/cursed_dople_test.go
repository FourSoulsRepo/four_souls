package r

import "testing"

func TestCursedDople(t *testing.T) {
	tb := fight(t, "cursed_dople")
	tb.Attack("cursed_dople", 6, 6, 6, 3)
	if tb.G.Players[1].Damage != 2 || tb.G.Players[2].Damage != 2 {
		t.Errorf("each damage deals 1 to each non-active player: %d %d", tb.G.Players[1].Damage, tb.G.Players[2].Damage)
	}
}
