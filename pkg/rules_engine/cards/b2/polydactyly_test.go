package b2

import "testing"

func TestPolydactyly(t *testing.T) {
	tb := itemTable(t, items("polydactyly"), "a_penny", "a_penny")
	tb.Play(0, "a_penny")
	tb.Play(0, "a_penny") // the additional loot play
	if c := tb.G.Players[0].Cents; c != 2 {
		t.Errorf("cents = %d, want 2", c)
	}
}
