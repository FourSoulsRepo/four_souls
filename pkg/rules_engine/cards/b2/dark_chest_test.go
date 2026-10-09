package b2

import "testing"

func TestDarkChest(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6)
	tb.RevealFromDeck("dark_chest")
	if !tb.G.Players[0].Dead {
		t.Error("2 damage did not kill Isaac")
	}
}
