package r

import "testing"

func TestMother(t *testing.T) {
	tb := fight(t, "mother")
	hurt(tb, 0, 2)
	tb.Attack("mother", 1, 2) // DC 1, then 2
	if len(tb.G.Monsters) != 4 {
		t.Errorf("slots = %d, want 2 + 2", len(tb.G.Monsters))
	}
}
