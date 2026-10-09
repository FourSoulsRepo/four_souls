package b2

import "testing"

func TestCursedChest(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6)
	tb.RevealFromDeck("cursed_chest", "guppys_head")
	if !hasItem(tb.G, 0, "guppys_head") {
		t.Error("no Guppy item gained")
	}
}
