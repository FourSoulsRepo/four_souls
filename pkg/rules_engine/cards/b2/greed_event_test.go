package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestGreedEventLostCentsAmount: the loss has its event, with the amount
// (19.192). The choice and the cents are in
// greed_event/richest-player-loses-all-cents.json.
func TestGreedEventLostCentsAmount(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.G.Players[1].Cents = 9
	events := tb.RevealFromDeck("greed_event", foe)
	lost := func(e engine.Event) bool { return e.Kind == engine.EvLostCents && e.Player == 1 && e.Amount == 9 }
	if !slices.ContainsFunc(events, lost) {
		t.Errorf("no lost_cents of 9¢ for Cain: %+v", events)
	}
}
