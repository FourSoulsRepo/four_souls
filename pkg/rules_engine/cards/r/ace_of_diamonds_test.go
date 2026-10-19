package r

import "testing"

func TestAceOfDiamonds(t *testing.T) {
	tb := lootTable(t, "ace_of_diamonds", "bomb", "a_nickel")
	tb.Play(0, "ace_of_diamonds")
	if got := tb.G.Players[0].Cents; got != 4 {
		t.Errorf("cents = %d, want 2 slots + 2 cards", got)
	}
}
