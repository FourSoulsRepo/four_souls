package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDeadWeight(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("the_fettered", "dead_weight"), seat("cain")},
		Monsters: []engine.CardRef{"leech", "fly"},
	})
	tb.Attack("leech", 1, 6)
	if tb.G.Players[0].Damage != 0 {
		t.Error("no combat damage on attack rolls of 1")
	}
}
