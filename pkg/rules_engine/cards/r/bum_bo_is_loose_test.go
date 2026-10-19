package r

import "testing"

func TestBumBoIsLoose(t *testing.T) {
	tb := roomTable(t, "bum_bo_is_loose")
	tb.G.Players[0].Cents = 5
	tb.EndTurn()
	if got := tb.G.Players[0].Cents; got != 3 {
		t.Errorf("cents = %d, want 3", got)
	}
}
