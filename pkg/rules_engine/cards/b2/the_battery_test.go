package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestTheBattery(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("the_battery", "the_d6"), Deactivated: items("the_d6")}, seat("cain"),
	}}, Set)
	tb.Start(0, "the_battery", 0)
	if opts := tb.Options(); len(opts) != 2 || opts[0] != "the_d6" {
		t.Errorf("options %v, want the D6 and cancel (not itself)", opts)
	}
	tb.Settle("the_d6")
	if !tb.G.Object(tb.Find(0, "the_d6")).Charged {
		t.Error("the D6 is not recharged")
	}
}
