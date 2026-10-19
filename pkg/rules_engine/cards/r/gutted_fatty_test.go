package r

import "testing"

func TestGuttedFatty(t *testing.T) {
	tb := fight(t, "gutted_fatty")
	hurt(tb, 0, 1)
	tb.Attack("gutted_fatty", 6)
	if len(tb.G.Monsters) != 3 || tb.G.Turn.MustAttacks != 1 {
		t.Errorf("slots %d, must attacks %d; want 3 and 1", len(tb.G.Monsters), tb.G.Turn.MustAttacks)
	}
}
