package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestClog(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			{Character: "isaac", Hand: []engine.CardRef{"bomb"}}, {Character: "cain", Hand: []engine.CardRef{"a_nickel"}},
		},
		Monsters: []engine.CardRef{"clog"},
	})
	hurt(tb, 0, 1)
	tb.Attack("clog", 6)
	if !handHas(tb.G, 1, "bomb") || !handHas(tb.G, 0, "a_nickel") {
		t.Error("the hands move left, then Isaac loots 3")
	}
	if len(tb.G.Players[0].Hand) != 4 {
		t.Errorf("Isaac's hand = %d, want the nickel + 3", len(tb.G.Players[0].Hand))
	}
}
