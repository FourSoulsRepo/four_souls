package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestConsolationPrize(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Souls: []engine.CardRef{"monstro", "gurdy"}}, seat("isaac", "consolation_prize"),
	}})
	tb.EndTurn()
	if got := len(tb.G.Players[1].Hand); got != 3 {
		t.Errorf("hand = %d, want 2 behind + 1 loot step", got)
	}
}
