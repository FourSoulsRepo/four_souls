package b2

import "testing"

func TestLostSoul(t *testing.T) {
	tb := lootTable(t, "lost_soul")
	tb.Play(0, "lost_soul")
	if s := tb.G.SoulValue(0); s != 1 {
		t.Errorf("soul value = %d, want 1", s)
	}
}
