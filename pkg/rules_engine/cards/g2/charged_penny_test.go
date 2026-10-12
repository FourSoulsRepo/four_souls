package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestChargedPenny(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"the_d6"}, Deactivated: []engine.CardRef{"the_d6"}, Hand: []engine.CardRef{"charged_penny"}},
		seat("cain"),
	}})
	tb.Play(0, "charged_penny", "the_d6")
	if c := tb.G.Players[0].Cents; c != 1 {
		t.Errorf("cents = %d, want 1", c)
	}
	if !tb.G.Object(tb.Find(0, "the_d6")).Charged {
		t.Error("the D6 was not recharged")
	}
}
