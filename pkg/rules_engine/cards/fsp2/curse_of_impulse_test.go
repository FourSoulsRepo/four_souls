package fsp2

import "testing"

func TestCurseOfImpulse(t *testing.T) {
	tb := armed(t, items("mystery_sack"), "fatty", "monstro")
	tb.RevealFromDeck("curse_of_impulse", me)
	tb.EndTurn()
	if tb.G.Object(tb.Find(0, "mystery_sack")).Charged || tb.G.Object(tb.G.Players[0].Character).Charged {
		t.Error("the items and the character are deactivated at the end of the turn")
	}
}
