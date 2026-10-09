package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestDeadBird(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "dead_bird"), {Character: "cain", Items: items("mystery_sack"), Hand: items("bomb")},
	}, Active: 1}, Set)
	tb.G.ForceRolls(3)
	tb.Activate(1, "mystery_sack", 0, "bomb")
	if h := handCards(tb.G, 0); len(h) != 1 || h[0] != "bomb" {
		t.Errorf("hand %v, want the stolen bomb", h)
	}
}
