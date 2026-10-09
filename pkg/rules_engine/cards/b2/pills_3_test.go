package b2

import "testing"

func TestPills3(t *testing.T) {
	tb := lootTable(t, "pills_3")
	tb.G.ForceRolls(1)
	tb.Play(0, "pills_3")
	if got := tb.G.PlayerATK(0); got != 2 {
		t.Errorf("roll 1: ATK %d, want 2", got)
	}
	tb = lootTable(t, "pills_3")
	tb.G.ForceRolls(3)
	tb.Play(0, "pills_3")
	if got := tb.G.PlayerHP(0); got != 3 {
		t.Errorf("roll 3: HP %d, want 3", got)
	}
	tb = lootTable(t, "pills_3")
	tb.G.ForceRolls(6)
	tb.Play(0, "pills_3")
	if got := tb.G.Players[0].Damage; got != 1 {
		t.Errorf("roll 6: damage %d, want 1", got)
	}
}
