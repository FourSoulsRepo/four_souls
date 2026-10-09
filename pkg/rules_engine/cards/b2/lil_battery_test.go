package b2

import (
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
