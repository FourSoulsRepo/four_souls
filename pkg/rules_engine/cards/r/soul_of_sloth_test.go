package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSoulOfSloth(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:    []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"the_d6", "jawbone", "spoon_bender"}, Cents: 20}, seat("cain")},
		BonusSouls: []engine.CardRef{"soul_of_sloth"},
	})
	tb.G.Turn.Purchases = 2
	// Isaac: 3 items and Book of Sin make 4; Cain has the fewest.
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "treasure deck", foe)
	settleAny(tb)
	if s := tb.G.SoulValue(1); s != 1 {
		t.Errorf("Cain's soul value = %d, want 1", s)
	}
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "treasure deck")
	settleAny(tb)
	if s := tb.G.SoulValue(1); s != 1 {
		t.Errorf("it happens only once: Cain's soul value = %d", s)
	}
}
