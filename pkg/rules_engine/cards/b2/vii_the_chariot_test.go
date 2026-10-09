package b2

import "testing"

func TestVIITheChariot(t *testing.T) {
	tb := lootTable(t, "vii_the_chariot")
	tb.Play(0, "vii_the_chariot", foe)
	if a, h := tb.G.PlayerATK(1), tb.G.PlayerHP(1); a != 2 || h != 3 {
		t.Errorf("ATK %d HP %d, want 2 and 3", a, h)
	}
}
