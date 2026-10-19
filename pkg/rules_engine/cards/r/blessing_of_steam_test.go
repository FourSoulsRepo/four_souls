package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBlessingOfSteam(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Cents: 5}, seat("cain")},
		Shop:    []engine.CardRef{"jawbone", "spoon_bender"},
		Rooms:   []engine.CardRef{"blessing_of_steam"},
	})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "jawbone")
	if !hasItem(tb.G, 0, "jawbone") || tb.G.Players[0].Cents != 0 {
		t.Errorf("the 10¢ jawbone costs 5¢: cents %d", tb.G.Players[0].Cents)
	}
}
