package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMoreOptions(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"more_options"}, Cents: 20}, seat("cain")},
		Shop:    []engine.CardRef{"the_d6", "breakfast"},
	})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "the_d6")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "breakfast")
	if c := tb.G.Players[0].Cents; c != 0 {
		t.Errorf("cents = %d, want 0 after two purchases", c)
	}
}
