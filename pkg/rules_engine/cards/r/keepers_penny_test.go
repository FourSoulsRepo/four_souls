package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestKeepersPenny(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Cents: 10}, seat("cain")}, Shop: []engine.CardRef{"keepers_penny"}})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "keepers_penny")
	if got := tb.G.Players[0].Cents; got != 9 {
		t.Fatalf("cents = %d, want 10 - 10 + 9", got)
	}
	tb.Activate(0, "keepers_penny", 1, foe)
	if got := tb.G.Players[0].Cents; got != 10 {
		t.Errorf("cents = %d, want 9 + 1 counter", got)
	}
	if !hasItem(tb.G, 1, "keepers_penny") || tb.G.Players[1].Cents != 0 {
		t.Error("given to Cain, who gains nothing: giving is not gaining")
	}
}
