package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestForeverAloneLooksAtADeck: the look is private and names the card
// (event fields). The other modes are situation files in forever_alone/.
func TestForeverAloneLooksAtADeck(t *testing.T) {
	tb := table(t, seat("blue_baby", "forever_alone"), seat("cain"))
	top := tb.G.Object(tb.G.DeckTop(engine.TreasureDeck, 1)[0]).Card
	ev := tb.Activate(0, "forever_alone", 0, "Look at the top card of a deck.", "ask.deck.treasure")
	saw := false
	for _, e := range ev {
		if e.Kind == engine.EvLookedAt && e.Player == 0 && e.Card == top && e.Private {
			saw = true
		}
	}
	if !saw {
		t.Errorf("no private look at %s", top)
	}
}
