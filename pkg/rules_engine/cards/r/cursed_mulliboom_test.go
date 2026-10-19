package r

import "testing"

func TestCursedMulliboom(t *testing.T) {
	tb := fight(t, "cursed_mulliboom")
	tb.Attack("fly", 6)
	sum := tb.G.Players[0].Damage + tb.G.Players[1].Damage + tb.G.Players[2].Damage
	if sum == 0 {
		t.Error("a monster died: 3 damage to a random player")
	}
}
