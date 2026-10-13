package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGuppysEye(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("guppys_eye"), Hand: items("bomb")},
		{Character: "cain", Hand: items("two_cents")},
		{Character: "maggy", Hand: items("four_cents")},
	}})
	v := tb.G.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 2})
	if len(v.Players[1].Hand) != 1 {
		t.Error("Cain's hand is revealed to everyone")
	}
	if len(v.Players[0].Hand) != 0 {
		t.Error("the Guppy's Eye player's own hand stays hidden")
	}
}
