package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestSteamySale(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: items("steamy_sale"), Cents: 5}, seat("cain")},
		Shop:    []engine.CardRef{"breakfast"},
	}, Set)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "breakfast")
	if !hasItem(tb.G, 0, "breakfast") || tb.G.Players[0].Cents != 0 {
		t.Error("could not buy for 5¢")
	}
}
