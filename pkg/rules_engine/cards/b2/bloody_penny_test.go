package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestBloodyPenny(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"bloody_penny"}, Hand: []engine.CardRef{"xiii_death"}},
		{Character: "cain", Cents: 3},
	}}, Set)
	tb.Play(0, "xiii_death", foe)
	// The loot trigger resolves before the penalty (R-DEATH-13).
	if n := len(tb.G.Players[0].Hand); n != 1 {
		t.Errorf("hand has %d cards, want 1", n)
	}
	if c := tb.G.Players[1].Cents; c != 2 {
		t.Errorf("the dead player has %d¢, want 2 after the penalty", c)
	}
}
