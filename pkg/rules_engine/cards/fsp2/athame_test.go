package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAthame(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("athame"), Hand: items("bomb")},
		{Character: "cain", Hand: items("a_penny", "two_cents", "four_cents")},
	}})
	tb.Activate(0, "athame", 0, "bomb", foe, "two_cents", "four_cents")
	if n := len(tb.G.Players[1].Hand); n != 1 {
		t.Errorf("Cain's hand = %d, want 1", n)
	}
}
