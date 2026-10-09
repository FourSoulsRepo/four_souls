package b2

import "testing"

func TestIITheHighPriestess(t *testing.T) {
	tb := lootTable(t, "ii_the_high_priestess")
	tb.G.ForceRolls(1)
	tb.Play(0, "ii_the_high_priestess", foe)
	if got := tb.G.Players[1].Damage; got != 1 {
		t.Errorf("damage = %d, want 1 (the roll)", got)
	}
}
