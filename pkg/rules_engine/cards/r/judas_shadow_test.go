package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestJudasShadow(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"mega_bomb"}}, seat("isaac", "judas_shadow"),
	}, Shop: []engine.CardRef{"spoon_bender"}})
	tb.Play(0, "mega_bomb", "player 2 (isaac)", "judas_shadow")
	if got := tb.G.Players[1].Cents; got != 2 {
		t.Errorf("2 counters, then destroyed as the penalty: cents %d, want 2", got)
	}
	if got := len(tb.G.Players[1].Hand); got != 2 {
		t.Errorf("hand = %d, want 2", got)
	}
}
