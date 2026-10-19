package r

import "testing"

func TestBlessingOfGluttony(t *testing.T) {
	tb := roomTable(t, "blessing_of_gluttony")
	n := len(tb.G.Players[1].Hand)
	tb.EndTurn()
	if got := len(tb.G.Players[1].Hand) - n; got != 2 {
		t.Errorf("Cain looted %d in the loot step, want 2", got)
	}
}
