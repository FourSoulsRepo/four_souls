package r

import "testing"

func TestCursedMulligan(t *testing.T) {
	tb := fight(t, "cursed_mulligan")
	tb.Attack("fly", 6)
	if len(tb.G.Monsters) != 3 {
		t.Errorf("slots = %d, want 3", len(tb.G.Monsters))
	}
}
