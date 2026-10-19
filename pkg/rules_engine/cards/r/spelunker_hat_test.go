package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSpelunkerHat(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("isaac", "spelunker_hat")}})
	tb.G.PutOnDeck(engine.TreasureDeck, "spoon_bender", "the_battery", "boomerang")
	tb.EndTurn("treasure deck", "spoon_bender", "the_battery", "boomerang")
	if top := tb.G.DeckTop(engine.TreasureDeck, 1); tb.G.Object(top[0]).Card != "spoon_bender" {
		t.Error("the Spoon Bender is on top now")
	}
}
