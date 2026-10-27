package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDamocles(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"damocles"}, Cents: 10}, seat("cain")}})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "ask.deck.treasure")
	if got := len(tb.G.Players[0].InPlay); got != 3 {
		t.Errorf("items = %d, want Damocles + 2 treasures", got)
	}
}

func TestDamoclesDestroyedOnDeath(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"mega_bomb"}}, seat("isaac", "damocles"),
	}})
	tb.Play(0, "mega_bomb", "player 2 (isaac)")
	if hasItem(tb.G, 1, "damocles") {
		t.Error("when he dies it is destroyed")
	}
}
