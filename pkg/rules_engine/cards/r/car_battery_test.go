package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCarBattery(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"car_battery", "the_d6", "sack_of_pennies"}, Deactivated: []engine.CardRef{"the_d6", "sack_of_pennies"}, Cents: 7},
		{Character: "cain", Items: []engine.CardRef{"jawbone"}, Deactivated: []engine.CardRef{"jawbone"}},
	}})
	tb.Activate(0, "car_battery", 0, "the_d6", "jawbone", "ask.done")
	if !tb.G.Object(tb.Find(0, "the_d6")).Charged || !tb.G.Object(tb.Find(1, "jawbone")).Charged || tb.G.Object(tb.Find(0, "sack_of_pennies")).Charged {
		t.Error("the D6 and the Jawbone are recharged, not the sack")
	}
	battery := tb.Find(0, "car_battery")
	tb.G.Recharge(battery)
	if tb.G.Object(battery).Charged {
		t.Error("only its own abilities recharge it")
	}
	tb.Activate(0, "car_battery", 1)
	if !tb.G.Object(battery).Charged || tb.G.Players[0].Cents != 0 {
		t.Error("$ Pay 7¢: recharge this")
	}
}

func TestCarBatteryRechargeStep(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("cain"), {Character: "isaac", Items: []engine.CardRef{"car_battery"}, Deactivated: []engine.CardRef{"car_battery"}},
	}})
	tb.EndTurn()
	if tb.G.Object(tb.Find(1, "car_battery")).Charged {
		t.Error("the recharge step doesn't recharge it")
	}
}
