package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestCounterfeitPenny(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"counterfeit_penny"}, Hand: []engine.CardRef{"a_nickel"}},
		seat("cain"),
	}}, Set)
	tb.Play(0, "a_nickel")
	if c := tb.G.Players[0].Cents; c != 6 {
		t.Errorf("cents = %d, want 6", c)
	}
}
