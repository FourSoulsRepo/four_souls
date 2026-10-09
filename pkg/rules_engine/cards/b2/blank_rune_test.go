package b2

import "testing"

func TestBlankRune(t *testing.T) {
	tb := lootTable(t, "blank_rune")
	tb.G.ForceRolls(6)
	tb.Play(0, "blank_rune")
	if a, b := tb.G.Players[0].Cents, tb.G.Players[1].Cents; a != 6 || b != 6 {
		t.Errorf("cents %d and %d, want 6 each", a, b)
	}
	tb = lootTable(t, "blank_rune")
	tb.G.ForceRolls(3)
	tb.Play(0, "blank_rune")
	if !tb.G.Players[0].Dead || !tb.G.Players[1].Dead {
		t.Error("3 damage to each player did not kill both")
	}
}
