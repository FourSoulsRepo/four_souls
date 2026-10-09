package b2

import "testing"

func TestStarterDeck(t *testing.T) {
	hand := items("a_dime", "a_dime", "a_dime", "a_dime", "a_dime", "a_dime", "a_dime", "a_dime")
	tb := itemTable(t, items("starter_deck"), hand...)
	tb.EndTurn()
	// Loot 2 at the end of the turn, then discard down to 10.
	if h := len(tb.G.Players[0].Hand); h != 10 {
		t.Errorf("hand = %d, want 10", h)
	}
}
