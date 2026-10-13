package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDaddyLongLegsFsp2(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), {Character: "cain", Items: items("daddy_long_legs_fsp2")}},
		Monsters: items("fly", "fatty"),
	})
	tb.EndTurn("fly") // Cain's turn starts: 1 damage to the fly
	if top(tb, 0) == "fly" {
		t.Error("the fly should die")
	}
}
