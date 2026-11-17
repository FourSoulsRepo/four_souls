package fsp2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

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

// Each deactivation has its event (19.176).
func TestCurseOfImpulseDeactivatedEvents(t *testing.T) {
	tb := armed(t, items("mystery_sack"), "fatty", "monstro")
	tb.RevealFromDeck("curse_of_impulse", me)
	char, sack := tb.G.Players[0].Character, tb.Find(0, "mystery_sack")
	ev := tb.EndTurn()
	for _, id := range []engine.ObjectID{char, sack} {
		if !slices.ContainsFunc(ev, func(e engine.Event) bool { return e.Kind == engine.EvDeactivated && e.Object == id }) {
			t.Errorf("no deactivated event names object %d", id)
		}
	}
}
