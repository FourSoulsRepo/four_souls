package b2

import "testing"

func TestGoatHead(t *testing.T) {
	tb := itemTable(t, items("goat_head"), "a_dime", "bomb", "a_nickel")
	tb.EndTurn("a_dime", "bomb", "done")
	if h := len(tb.G.Players[0].Hand); h != 3 {
		t.Errorf("hand %d, want 3", h)
	}
	if !inDiscard(tb.G, "a_dime") || !inDiscard(tb.G, "bomb") {
		t.Error("the chosen cards are not discarded")
	}
}
