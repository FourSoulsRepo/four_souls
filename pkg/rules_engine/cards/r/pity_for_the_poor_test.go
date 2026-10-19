package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPityForThePoor(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Shop:    []engine.CardRef{"jawbone", "spoon_bender"},
		Rooms:   []engine.CardRef{"pity_for_the_poor"},
	})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "jawbone")
	if !hasItem(tb.G, 0, "jawbone") {
		t.Error("Isaac (fewest souls, 0¢) buys the jawbone for 0¢ (Q-21)")
	}
}
