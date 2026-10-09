package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestVIIIJustice(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"viii_justice"}, Cents: 1},
		{Character: "cain", Hand: []engine.CardRef{"a_dime", "a_dime", "bomb"}, Cents: 9},
	}}, Set)
	tb.Play(0, "viii_justice", foe)
	if h, c := len(tb.G.Players[0].Hand), tb.G.Players[0].Cents; h != 3 || c != 9 {
		t.Errorf("hand %d cents %d, want 3 and 9", h, c)
	}
}
