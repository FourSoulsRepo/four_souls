package r

import "testing"

func TestBlessingOfTheSack(t *testing.T) {
	tb := roomTable(t, "blessing_of_the_sack")
	tb.Attack("fly", 6)
	if got := tb.G.Players[0].Cents; got != 2 {
		t.Errorf("cents = %d, want 2 (the fly's 1¢ doubled)", got)
	}
}
