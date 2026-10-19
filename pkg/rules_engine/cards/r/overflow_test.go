package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestOverflow(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "overflow")
	tb.Attack("fly", 6, 3, 1) // Isaac rolls 3, Cain rolls 1
	if tb.G.Players[1].SkipTurns != 1 || tb.G.Players[0].SkipTurns != 0 {
		t.Error("Cain skips his next turn")
	}
}
