package b2

import "testing"

func TestEdensBlessing(t *testing.T) {
	tb := itemTable(t, items("edens_blessing"))
	tb.EndTurn()
	if c := tb.G.Players[0].Cents; c != 6 {
		t.Errorf("cents = %d, want 6", c)
	}
	tb = itemTable(t, items("edens_blessing"))
	tb.G.Players[0].Cents = 1
	tb.EndTurn()
	if c := tb.G.Players[0].Cents; c != 1 {
		t.Errorf("with 1¢: cents = %d, want 1", c)
	}
}
