package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSumptorium(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("the_curdled", "sumptorium"), seat("cain")},
		Monsters: []engine.CardRef{"leech", "fly"},
	})
	s := tb.Find(0, "sumptorium")
	tb.G.AddCountersTo(s, 2)
	tb.Activate(0, "sumptorium", 1, "leech")
	if tb.G.Evasion(tb.G.Monsters[0].Cards[0]) != 3 {
		t.Error("the leech has -1 DC")
	}
}
