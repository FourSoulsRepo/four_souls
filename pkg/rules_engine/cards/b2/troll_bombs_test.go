package b2

import "testing"

func TestTrollBombs(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.RevealFromDeck("troll_bombs")
	if !tb.G.Players[0].Dead || tb.G.Players[1].Dead {
		t.Error("only the active player takes 2 damage")
	}
}
