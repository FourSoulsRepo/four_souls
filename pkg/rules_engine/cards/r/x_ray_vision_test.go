package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestXRayVision(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "x_ray_vision"), {Character: "cain", Hand: []engine.CardRef{"bomb"}}}})
	tb.G.PutOnDeck(engine.TreasureDeck, "spoon_bender")
	ev := tb.Activate(0, "x_ray_vision", 0, foe, "treasure deck", "put spoon_bender into discard")
	saw := false
	for _, e := range ev {
		if e.Kind == engine.EvLookedAt && e.Player == 0 && e.Card == "bomb" {
			saw = true
		}
	}
	if !saw {
		t.Error("Isaac sees Cain's hand")
	}
	if !inDiscard(tb.G, engine.TreasureDeck, "spoon_bender") {
		t.Error("the top card goes to discard")
	}
	tb.EndTurn()
	if !tb.G.Object(tb.Find(0, "x_ray_vision")).Charged {
		t.Error("it recharges at the end of his turn")
	}
}
