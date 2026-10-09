package b2

import "testing"

func TestDevilDeal(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.RevealFromDeck("devil_deal", "loot 2, take 1 damage")
	if h, d := len(tb.G.Players[0].Hand), tb.G.Players[0].Damage; h != 2 || d != 1 {
		t.Errorf("hand %d damage %d, want 2 and 1", h, d)
	}
}
