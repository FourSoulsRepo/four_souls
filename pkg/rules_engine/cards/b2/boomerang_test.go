package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestBoomerang(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "boomerang"), {Character: "cain", Hand: items("a_dime", "bomb")},
	}}, Set)
	tb.Activate(0, "boomerang", 0, foe)
	if a, b := len(tb.G.Players[0].Hand), len(tb.G.Players[1].Hand); a != 1 || b != 1 {
		t.Errorf("hands %d and %d, want 1 and 1", a, b)
	}
}
