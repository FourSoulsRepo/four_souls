package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestDarkBum(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac"), seat("cain", "dark_bum")}}, Set)
	tb.G.ForceRolls(1)
	tb.EndTurn()
	if c := tb.G.Players[1].Cents; c != 3 {
		t.Errorf("cents = %d, want 3", c)
	}
}
