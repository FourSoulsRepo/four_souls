package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestFetalHaunt(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"fetal_haunt"}, Hand: []engine.CardRef{"bomb"}}, seat("cain"),
	}})
	tb.EndTurn("bomb", foe)
	if !handHas(tb.G, 1, "bomb") {
		t.Error("Isaac gives Cain the bomb")
	}
}
