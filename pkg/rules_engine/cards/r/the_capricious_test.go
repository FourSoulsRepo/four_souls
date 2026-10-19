package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheCapriciousCharacter(t *testing.T) {
	testCharacter(t, "the_capricious", "glitch", "Play an additional loot card this turn.")
}

func TestTheCapriciousShopItem(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("the_capricious"), seat("cain")}, Shop: []engine.CardRef{"spoon_bender"}})
	tb.Activate(0, "the_capricious", 0, "Put a shop item into discard.", "spoon_bender")
	if !inDiscard(tb.G, engine.TreasureDeck, "spoon_bender") {
		t.Error("the Spoon Bender goes to discard")
	}
}
