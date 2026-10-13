package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAAABattery(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"aaa_battery", "mystery_sack"}, Deactivated: []engine.CardRef{"mystery_sack"}}, seat("cain"),
	}})
	tb.G.ForceRolls(5)
	tb.EndTurn("mystery_sack")
	if !tb.G.Object(tb.Find(0, "mystery_sack")).Charged {
		t.Error("4-6: the sack should be recharged")
	}
}
