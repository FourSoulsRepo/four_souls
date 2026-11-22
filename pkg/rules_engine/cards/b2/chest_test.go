package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestChestDiscardedEvent: the done event card's discarded event names
// the card and its object in the slot (19.173), so the client moves that
// card to the discard. The cents and the uncovered Fly are in
// chest/gains-6-cents-on-a-5.json.
func TestChestDiscardedEvent(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.G.ForceRolls(5)
	ev := tb.RevealFromDeck("chest")
	if !slices.ContainsFunc(ev, func(e engine.Event) bool {
		return e.Kind == engine.EvDiscarded && e.Card == "chest" && e.Prev != 0
	}) {
		t.Error("no discarded event for the done event card (19.173)")
	}
}
