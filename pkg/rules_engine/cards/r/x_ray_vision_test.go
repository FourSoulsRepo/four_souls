package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The rest is x_ray_vision/discards-a-deck-top-and-recharges.json; the
// looked_at event's Card stays Go: no check reads event fields (24.6).
func TestXRayVisionShowsTheHand(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "x_ray_vision"), {Character: "cain", Hand: []engine.CardRef{"bomb"}}}})
	tb.G.PutOnDeck(engine.TreasureDeck, "spoon_bender")
	ev := tb.Activate(0, "x_ray_vision", 0, foe, "ask.deck.treasure", "ask.put_into_discard spoon_bender")
	saw := false
	for _, e := range ev {
		if e.Kind == engine.EvLookedAt && e.Player == 0 && e.Card == "bomb" {
			saw = true
		}
	}
	if !saw {
		t.Error("Isaac sees Cain's hand")
	}
}
