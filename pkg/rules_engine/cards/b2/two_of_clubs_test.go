package b2

import "testing"

func TestTwoOfClubs(t *testing.T) {
	tb := itemTable(t, items("two_of_clubs", "bum_friend"), "a_dime")
	tb.Activate(0, "two_of_clubs", 0, me)
	tb.Activate(0, "bum_friend", 0, "a_dime") // loot 1 doubled, put one back
	if h := len(tb.G.Players[0].Hand); h != 2 {
		t.Errorf("hand = %d, want 2", h)
	}
}
