package b2

import "testing"

func TestDagaz(t *testing.T) {
	tb := lootTable(t, "dagaz")
	tb.Play(0, "dagaz", "Choose a player. Prevent the next 1 damage they would take this turn.", me)
	tb.Attack("fly", 1, 6)
	if d := tb.G.Players[0].Damage; d != 0 {
		t.Errorf("damage = %d, want 0", d)
	}
}
