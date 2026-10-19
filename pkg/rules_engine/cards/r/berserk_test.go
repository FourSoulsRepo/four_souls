package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBerserk(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("the_savage", "berserk"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.Attack("fly", 6)
	if tb.G.PlayerHP(0) != 3 {
		t.Errorf("a kill: +1 HP till end of turn; HP %d", tb.G.PlayerHP(0))
	}
	tb.Attack("leech", 6)
	if top(tb, 1) == "leech" {
		t.Error("an additional attack on his turn")
	}
}
