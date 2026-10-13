package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCreditCard(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"credit_card"}, Cents: 3}, seat("cain")},
		Shop:    []engine.CardRef{"the_d6", "breakfast"},
	})
	tb.Play(0, "credit_card")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "the_d6")
	if !hasItem(tb.G, 0, "the_d6") {
		t.Fatal("the free purchase failed")
	}
	if c := tb.G.Players[0].Cents; c != 3 {
		t.Errorf("cents = %d, want 3: the first shop item is free", c)
	}
}

func TestCreditCardAfterAPurchase(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"theres_options"}, Hand: []engine.CardRef{"credit_card"}, Cents: 13}, seat("cain")},
		Shop:    []engine.CardRef{"the_d6", "breakfast"},
	})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "the_d6")
	tb.Play(0, "credit_card")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "breakfast")
	if c := tb.G.Players[0].Cents; c != 3 {
		t.Errorf("cents = %d, want 3: the second item still costs 10¢", c)
	}
}
