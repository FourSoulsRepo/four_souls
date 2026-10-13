package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMomsHeart(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), {Character: "cain", Souls: items("gurdy", "monstro")}},
		Monsters: items("moms_heart", "fly"),
	})
	tb.G.KillObject(topID(tb, 0))
	tb.Settle("player 1 (isaac)") // Isaac gains 2 souls: tied with Cain; he picks himself
	if !tb.G.Over || len(tb.G.Winners) != 1 || tb.G.Winners[0] != 0 {
		t.Errorf("over %v winners %v, want Isaac", tb.G.Over, tb.G.Winners)
	}
}
