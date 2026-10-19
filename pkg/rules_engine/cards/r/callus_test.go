package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCallusStopsNonCombat(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"bomb"}}, {Character: "isaac", Items: []engine.CardRef{"callus"}},
	}})
	tb.Play(0, "bomb", "player 2 (isaac)")
	if got := tb.G.Players[1].Damage; got != 0 {
		t.Errorf("a bomb is not combat: damage %d, want none", got)
	}
}

func TestCallusAllowsCombat(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"callus"}}, seat("cain")},
		Monsters: []engine.CardRef{"fly"},
	})
	tb.Attack("fly", 1, 6) // the fly hits back for 1: combat
	if got := tb.G.Players[0].Damage; got != 1 {
		t.Errorf("combat damage = %d, want 1", got)
	}
}
