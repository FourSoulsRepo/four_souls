package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestForgetMeNow(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("forget_me_now"), Souls: items("monstro")},
		{Character: "cain", Souls: items("gurdy", "fatty")},
	}})
	tb.Activate(0, "forget_me_now", 0, "monstro", "gurdy")
	if tb.G.SoulValue(0) != 0 || tb.G.SoulValue(1) != 1 {
		t.Errorf("souls = %d, %d, want 0, 1", tb.G.SoulValue(0), tb.G.SoulValue(1))
	}
}
