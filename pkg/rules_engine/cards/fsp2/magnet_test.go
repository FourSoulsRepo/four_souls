package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMagnet(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: items("three_cents")}, {Character: "cain", Items: items("magnet")},
	}})
	tb.Play(0, "three_cents")
	if tb.G.Players[0].Cents != 2 || tb.G.Players[1].Cents != 1 {
		t.Errorf("cents = %d, %d, want 2, 1", tb.G.Players[0].Cents, tb.G.Players[1].Cents)
	}
}
