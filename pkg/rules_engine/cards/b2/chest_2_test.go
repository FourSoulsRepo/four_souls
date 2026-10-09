package b2

import "testing"

func TestChest2(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6)
	tb.RevealFromDeck("chest_2")
	if h := len(tb.G.Players[0].Hand); h != 3 {
		t.Errorf("hand = %d, want 3", h)
	}
}
