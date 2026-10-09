package b2

import "testing"

func TestPandorasBox(t *testing.T) {
	tb := itemTable(t, items("pandoras_box"))
	tb.G.ForceRolls(6)
	tb.Activate(0, "pandoras_box", 0)
	if s := tb.G.SoulValue(0); s != 1 {
		t.Errorf("soul value = %d, want 1", s)
	}
	tb = itemTable(t, items("pandoras_box"))
	tb.G.ForceRolls(5)
	tb.Activate(0, "pandoras_box", 0)
	if c := tb.G.Players[0].Cents; c != 9 || hasItem(tb.G, 0, "pandoras_box") {
		t.Errorf("cents %d", c)
	}
}
