package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestRKey(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"r_key"}, Souls: []engine.CardRef{"monstro"}, Hand: []engine.CardRef{"bomb"}},
		{Character: "cain", Souls: []engine.CardRef{"gurdy"}},
	}})
	tb.Activate(0, "r_key", 0)
	if tb.G.SoulValue(0) != 0 || tb.G.SoulValue(1) != 0 {
		t.Error("every soul is destroyed")
	}
	if len(tb.G.Players[0].Hand) != 3 || len(tb.G.Players[1].Hand) != 3 || handHas(tb.G, 0, "bomb") {
		t.Error("each player discards their hand and loots 3")
	}
}

func TestRKeyEntersDeactivated(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Cents: 10}, seat("cain")}, Shop: []engine.CardRef{"r_key"}})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "r_key")
	if tb.G.Object(tb.Find(0, "r_key")).Charged {
		t.Error("it enters play deactivated")
	}
}
