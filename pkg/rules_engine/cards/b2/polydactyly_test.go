package b2

import "testing"

func TestPolydactyly(t *testing.T) {
	tb := itemTable(t, items("polydactyly"), "a_penny_6", "a_penny_6")
	tb.Play(0, "a_penny_6")
	tb.Play(0, "a_penny_6") // the additional loot play
	if c := tb.G.Players[0].Cents; c != 2 {
		t.Errorf("cents = %d, want 2", c)
	}
}
