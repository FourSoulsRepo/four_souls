package r

import "testing"

func TestMushroom(t *testing.T) {
	tb := fight(t, "mushroom")
	tb.Attack("mushroom", 1, 6, 6, 6)
	if tb.G.Players[1].Damage != 1 || tb.G.Players[2].Damage != 1 {
		t.Error("its combat damage to Isaac deals 1 to each other player")
	}
}
