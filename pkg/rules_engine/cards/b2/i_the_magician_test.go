package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestITheMagician(t *testing.T) {
	tb := rollTable(t, 2, engine.SituationPlayer{Character: "isaac", Hand: []engine.CardRef{"i_the_magician"}}, seat("cain"))
	tb.Pass(1)
	if got := resolvedRoll(t, tb.Play(0, "i_the_magician", "roll of 2", "5")); got != 5 {
		t.Errorf("resolved as %d, want 5", got)
	}
}
