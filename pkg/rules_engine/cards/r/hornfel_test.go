package r

import "testing"

// The player Hornfel hurts is chosen at random, so no file can name it;
// hornfel/combat-damage-to-isaac-is-prevented.json checks the rest (24.24).
func TestHornfelHurtsOnePlayerAtRandom(t *testing.T) {
	tb := fight(t, "hornfel")
	tb.Attack("hornfel", 1, 6, 6)
	if tb.G.Players[0].Damage+tb.G.Players[1].Damage+tb.G.Players[2].Damage != 1 {
		t.Error("its combat damage is prevented; 1 damage to a random player instead")
	}
}
