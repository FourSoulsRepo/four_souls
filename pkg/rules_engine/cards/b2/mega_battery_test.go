package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestMegaBattery(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"mega_battery"}},
		{Character: "cain", Items: []engine.CardRef{"the_d6", "blood_lust"}, Deactivated: []engine.CardRef{"the_d6", "blood_lust"}},
	}}, Set)
	tb.Play(0, "mega_battery", foe)
	for _, c := range []engine.CardRef{"the_d6", "blood_lust"} {
		if !tb.G.Object(tb.Find(1, c)).Charged {
			t.Errorf("%s is not recharged", c)
		}
	}
}
