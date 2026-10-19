package r

import "testing"

func TestBlessingOfGreed(t *testing.T) {
	tb := roomTable(t, "blessing_of_greed")
	c := tb.G.Players[1].Cents
	tb.EndTurn()
	if got := tb.G.Players[1].Cents - c; got != 3 {
		t.Errorf("Cain gained %d¢ at the start of his turn, want 3", got)
	}
}
