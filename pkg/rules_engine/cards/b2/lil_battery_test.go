package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestLilBattery(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"the_d6"}, Deactivated: []engine.CardRef{"the_d6"}, Hand: []engine.CardRef{"lil_battery"}},
		seat("cain"),
	}}, Set)
	tb.Play(0, "lil_battery", "the_d6")
	if !tb.G.Object(tb.Find(0, "the_d6")).Charged {
		t.Error("the item is not recharged")
	}
}

func TestLilBatteryOffersOnlyActiveItems(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("the_d6", "breakfast", "battery_bum"), Hand: []engine.CardRef{"lil_battery"}},
		seat("cain"),
	}}, Set)
	// Apply, not Play: Play would fail on the open target prompt.
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "lil_battery")}}); err != nil {
		t.Fatal(err)
	}
	got := tb.Options()
	if !slices.Contains(got, "the_d6") {
		t.Errorf("options = %v: the active D6 is not offered", got)
	}
	for _, card := range []string{"breakfast", "battery_bum"} {
		if slices.Contains(got, card) {
			t.Errorf("options = %v: %s is offered, but it has no ↷ ability", got, card)
		}
	}
}
