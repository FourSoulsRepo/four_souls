package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTwoOfSpades(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"two_of_spades"}},
		{Character: "cain", Hand: []engine.CardRef{"bomb", "bomb", "a_nickel"}},
	}})
	tb.Play(0, "two_of_spades", foe)
	if got := len(tb.G.Players[0].Hand); got != 3 {
		t.Errorf("hand = %d, want 3 like Cain's", got)
	}
}
