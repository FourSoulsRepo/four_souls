package b2

import "testing"

func TestBlankCard(t *testing.T) {
	tb := itemTable(t, items("blank_card"), "a_nickel")
	tb.Activate(0, "blank_card", 0)
	tb.Play(0, "a_nickel")
	if c := tb.G.Players[0].Cents; c != 10 {
		t.Errorf("cents = %d, want 10 (copied)", c)
	}
}
