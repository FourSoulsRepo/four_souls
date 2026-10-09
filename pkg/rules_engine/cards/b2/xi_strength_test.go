package b2

import "testing"

func TestXIStrength(t *testing.T) {
	tb := lootTable(t, "xi_strength")
	tb.Play(0, "xi_strength", me)
	if a, n := tb.G.PlayerATK(0), tb.G.Turn.Attacks; a != 2 || n != 2 {
		t.Errorf("ATK %d attacks %d, want 2 and 2", a, n)
	}
}
