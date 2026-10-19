package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestKeepersSack(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Cents: 10}, seat("cain")}, Shop: []engine.CardRef{"keepers_sack"}})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "keepers_sack")
	if got := tb.G.Players[0].Cents; got != 14 {
		t.Fatalf("cents = %d, want 14", got)
	}
	tb.Activate(0, "keepers_sack", 1, foe)
	if !hasItem(tb.G, 1, "keepers_sack") || tb.G.Players[0].Cents != 10 || tb.G.Players[1].Cents != 0 {
		t.Error("Pay 4¢: given to Cain, who gains nothing (S-YUG)")
	}
}

func TestKeepersSackEndOfTurn(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"keepers_sack"}}, seat("cain")}})
	tb.EndTurn()
	if got := tb.G.Players[0].Damage; got != 0 {
		t.Errorf("damage %d after the turn: healed at its end", got)
	}
	tb2 := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"keepers_sack"}, Hand: []engine.CardRef{"bomb"}}, seat("cain")}})
	tb2.EndTurn("bomb")
	if len(tb2.G.Players[0].Hand) != 0 {
		t.Error("he discards the bomb")
	}
}
