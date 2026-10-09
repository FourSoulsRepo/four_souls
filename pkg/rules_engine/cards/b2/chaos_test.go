package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestChaos(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("chaos"), Hand: items("a_dime")}, {Character: "cain", Hand: items("bomb", "bomb")},
	}}, Set)
	tb.Activate(0, "chaos", 0)
	if a, b := handCards(tb.G, 0), handCards(tb.G, 1); len(a) != 2 || len(b) != 1 || b[0] != "a_dime" {
		t.Errorf("hands %v and %v", a, b)
	}
}
