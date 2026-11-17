package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestBatteryBum(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("battery_bum", "the_d6"), Deactivated: items("the_d6"), Cents: 4},
		seat("cain"),
	}}, Set)
	tb.Activate(0, "battery_bum", 0, "the_d6")
	if !tb.G.Object(tb.Find(0, "the_d6")).Charged || tb.G.Players[0].Cents != 0 {
		t.Error("4¢ did not recharge the D6")
	}
	if rule := tb.Refused(0, "battery_bum", 0); rule != "R-ABIL-07" {
		t.Errorf("without 4¢: refused by %s, want R-ABIL-07", rule)
	}
}
