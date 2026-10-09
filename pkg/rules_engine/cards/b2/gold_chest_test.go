package b2

import "testing"

func TestGoldChest(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.G.ForceRolls(1)
	tb.RevealFromDeck("gold_chest")
	if n := len(tb.G.Players[0].InPlay); n != 1 {
		t.Errorf("%d items, want 1", n)
	}
}
