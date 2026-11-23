package fsp2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Each deactivation has its event (19.176): the event's object is a
// field no situation check reads. The deactivated cards are in
// testdata/situations/curse_of_impulse/.
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
