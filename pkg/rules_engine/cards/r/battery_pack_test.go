package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBatteryPack(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"battery_pack", "the_d6", "spoon_bender"}, Deactivated: []engine.CardRef{"the_d6", "spoon_bender"}}, seat("cain"),
	}})
	if rule := tb.Refused(0, "battery_pack", 1); rule != "R-ABIL-07" {
		t.Errorf("no counters: refused by %s", rule)
	}
	tb.G.AddCountersTo(tb.Find(0, "battery_pack"), 2)
	tb.Activate(0, "battery_pack", 1, me)
	if !tb.G.Object(tb.Find(0, "the_d6")).Charged || !tb.G.Object(tb.Find(0, "spoon_bender")).Charged {
		t.Error("each of Isaac's items is recharged")
	}
}
