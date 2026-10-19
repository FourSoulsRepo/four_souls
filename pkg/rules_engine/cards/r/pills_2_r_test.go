package r

import "testing"

func TestPillsRollThreeTimes(t *testing.T) {
	tb := lootTable(t, "pills_2_r")
	tb.G.ForceRolls(1, 3, 5)
	tb.Play(0, "pills_2_r")
	if got := tb.G.PlayerATK(0); got != 2 {
		t.Errorf("ATK = %d, want 1 + 1", got)
	}
	if got := tb.G.Players[0].Damage; got != 1 {
		t.Errorf("damage = %d, want 1", got)
	}
	if got := tb.G.PlayerHP(0); got != 2 {
		t.Errorf("HP left = %d, want 2 + 1 - 1", got)
	}
}
