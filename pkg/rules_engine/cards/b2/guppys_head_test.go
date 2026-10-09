package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestGuppysHead(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "guppys_head"), {Character: "cain", Hand: items("a_dime", "bomb")},
	}}, Set)
	tb.Activate(0, "guppys_head", 0, foe, "bomb") // Cain picks the card
	if h := handCards(tb.G, 0); len(h) != 1 || h[0] != "bomb" {
		t.Errorf("hand %v, want the bomb", h)
	}
}
