package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMomsBottleOfPills(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Cents: 10}, seat("cain")}})
	tb.G.PutOnDeck(engine.LootDeck, "pills", "bomb", "pills_r")
	tb.G.PutOnDeck(engine.TreasureDeck, "moms_bottle_of_pills")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "ask.deck.treasure")
	if !handHas(tb.G, 0, "pills") || !handHas(tb.G, 0, "pills_r") || handHas(tb.G, 0, "bomb") {
		t.Error("each Pills! goes to his hand")
	}
}
