package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTwoOfDiamonds(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"two_of_diamonds"}, Cents: 2},
		{Character: "cain", Cents: 5},
	}})
	tb.Play(0, "two_of_diamonds", foe)
	if c := tb.G.Players[1].Cents; c != 10 {
		t.Errorf("Cain's cents = %d, want 10", c)
	}
	if c := tb.G.Players[0].Cents; c != 2 {
		t.Errorf("Isaac's cents = %d, want 2", c)
	}
}
