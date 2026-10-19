package r

import "testing"

func TestMegaBomb(t *testing.T) {
	tb := lootTable(t, "mega_bomb")
	tb.Play(0, "mega_bomb", foe)
	if tb.G.Players[1].Damage != 2 || !tb.G.Players[1].Dead {
		t.Errorf("Cain: damage %d, dead %v; want 2 damage and dead", tb.G.Players[1].Damage, tb.G.Players[1].Dead)
	}
}
