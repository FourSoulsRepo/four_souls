package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestRedGhost(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"bomb"}}, seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"red_ghost"},
	})
	tb.Play(0, "bomb", "red_ghost", "player 3 (samson)")
	if tb.G.Players[2].Damage != 1 {
		t.Error("the active player chose Samson: 1 damage")
	}
}
