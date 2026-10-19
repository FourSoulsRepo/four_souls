package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestKeepersBox(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "cain", Cents: 10}, seat("isaac", "keepers_box")},
		Shop:    []engine.CardRef{"spoon_bender"},
	})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "spoon_bender")
	if tb.G.Players[1].Cents != 2 || len(tb.G.Players[1].Hand) != 1 {
		t.Error("another player bought a shop item: 2¢ and loot 1")
	}
	tb2 := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Cents: 10}, seat("cain")}})
	tb2.G.PutOnDeck(engine.TreasureDeck, "keepers_box")
	shop := len(tb2.G.Shop)
	tb2.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "treasure deck")
	if len(tb2.G.Shop) != shop+2 {
		t.Error("it entered play: 2 more shop slots")
	}
}
