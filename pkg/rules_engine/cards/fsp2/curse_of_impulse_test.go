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

// A passive item has no ↷ ability, so it stays charged (19.64).
func TestCurseOfImpulseSkipsPassiveItems(t *testing.T) {
	tb := armed(t, items("mystery_sack", "dads_lost_coin"), "fatty", "monstro")
	tb.RevealFromDeck("curse_of_impulse", me)
	tb.EndTurn()
	if !tb.G.Object(tb.Find(0, "dads_lost_coin")).Charged {
		t.Error("Dad's Lost Coin is passive: it stays charged")
	}
}
