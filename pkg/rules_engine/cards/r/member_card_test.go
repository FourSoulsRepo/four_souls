package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMemberCard(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"member_card"}, Cents: 6}, seat("cain")},
		Shop:    []engine.CardRef{"spoon_bender", "sack_of_pennies"},
	})
	tb.G.PutOnDeck(engine.TreasureDeck, "the_battery", "boomerang", "mom_s_bra")
	tb.Activate(0, "member_card", 0, "the_battery", "boomerang", "mom_s_bra", "spoon_bender")
	if !inDiscard(tb.G, engine.TreasureDeck, "spoon_bender") {
		t.Error("the Spoon Bender should be in the treasure discard")
	}
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "sack_of_pennies")
	if !hasItem(tb.G, 0, "sack_of_pennies") || tb.G.Players[0].Cents != 0 {
		t.Errorf("bought for 6¢: cents left %d", tb.G.Players[0].Cents)
	}
}
