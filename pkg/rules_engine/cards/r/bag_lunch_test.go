package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBagLunch(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"berkano", "berkano"}}, seat("cain", "bag_lunch"),
	}, Shop: []engine.CardRef{"spoon_bender", "sack_of_pennies"}})
	tb.G.Turn.LootPlays = 2
	if got := tb.G.PlayerHP(1); got != 2 {
		t.Fatalf("HP = %d, want 2 without a counter", got)
	}
	tb.Play(0, "berkano", "bag_lunch", "spoon_bender")
	if !hasItem(tb.G, 1, "bag_lunch") {
		t.Fatal("the first destruction puts a counter on it instead")
	}
	if got := tb.G.PlayerHP(1); got != 3 {
		t.Errorf("HP = %d, want 2 + 1 with a counter", got)
	}
	tb.Play(0, "berkano", "bag_lunch", "sack_of_pennies")
	if hasItem(tb.G, 1, "bag_lunch") {
		t.Error("with a counter it is destroyed")
	}
}
