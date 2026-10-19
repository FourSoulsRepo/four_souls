package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTrickPenny(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"trick_penny"}, Cents: 9}, seat("cain")},
		Shop:    []engine.CardRef{"spoon_bender"},
	})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "spoon_bender")
	if !hasItem(tb.G, 0, "spoon_bender") || tb.G.Players[0].Cents != 1 {
		t.Errorf("bought for 9¢, then +1¢: cents %d", tb.G.Players[0].Cents)
	}
	tb.EndTurn()
	if tb.G.Players[0].Cents != 2 {
		t.Error("+1¢ at the end of his turn")
	}
}

func TestTrickPennyPaid(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"trick_penny", "the_jar"}, Cents: 3}, seat("cain"),
	}})
	tb.Activate(0, "the_jar", 0, foe)
	if tb.G.Players[0].Cents != 0 {
		t.Errorf("The Jar's 4¢ costs 3¢: cents %d", tb.G.Players[0].Cents)
	}
}
