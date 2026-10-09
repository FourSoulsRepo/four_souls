package b2

import "testing"

func TestMegaTrollBomb(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.RevealFromDeck("mega_troll_bomb")
	if !tb.G.Players[0].Dead || !tb.G.Players[1].Dead {
		t.Error("2 damage each did not kill both")
	}
}
