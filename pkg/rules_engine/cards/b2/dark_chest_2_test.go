package b2

import "testing"

func TestDarkChest2(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.G.ForceRolls(3)
	tb.RevealFromDeck("dark_chest_2")
	if h := len(tb.G.Players[0].Hand); h != 2 {
		t.Errorf("hand = %d, want 2", h)
	}
}
