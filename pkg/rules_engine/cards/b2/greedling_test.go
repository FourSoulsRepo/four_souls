package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGreedling(t *testing.T) {
	tb := slayTable(t, "greedling", seat("isaac"), engine.SituationPlayer{Character: "cain", Cents: 8})
	killWithSixes(t, tb, foe)
	if c := tb.G.Players[1].Cents; c != 1 {
		t.Errorf("Cain's cents = %d, want 1", c)
	}
}
