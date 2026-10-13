package fsp2

import "testing"

func TestPills3Fsp2(t *testing.T) {
	tb := lootTable(t, "pills_3_fsp2")
	tb.G.ForceRolls(3)
	tb.Play(0, "pills_3_fsp2", "fly")
	if top(tb, 0) == "fly" {
		t.Error("3-4: the fly should take 1 damage and die")
	}
	tb2 := lootTable(t, "pills_3_fsp2")
	tb2.G.ForceRolls(6)
	tb2.Play(0, "pills_3_fsp2")
	if d := tb2.G.Players[0].Damage; d != 1 {
		t.Errorf("5-6: damage = %d, want 1", d)
	}
}
