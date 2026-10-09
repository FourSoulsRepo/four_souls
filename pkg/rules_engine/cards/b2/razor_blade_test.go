package b2

import "testing"

func TestRazorBlade(t *testing.T) {
	tb := itemTable(t, items("razor_blade"))
	tb.Activate(0, "razor_blade", 0, foe)
	if d := tb.G.Players[1].Damage; d != 1 {
		t.Errorf("damage = %d, want 1", d)
	}
}
