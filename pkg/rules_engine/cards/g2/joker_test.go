package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestJoker(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"joker"}},
		{Character: "cain", Hand: []engine.CardRef{"bomb", "a_sack"}},
	}})
	tb.Play(0, "joker", foe, "a_sack")
	if len(tb.G.Players[0].Hand) != 1 || tb.G.Object(tb.G.Players[0].Hand[0]).Card != "a_sack" {
		t.Errorf("Isaac's hand = %v, want the stolen A Sack", tb.G.Players[0].Hand)
	}
	if len(tb.G.Players[1].Hand) != 1 {
		t.Error("Cain should have 1 card left")
	}
	tb2 := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"joker"}},
		{Character: "cain", Hand: []engine.CardRef{"bomb"}},
	}})
	tb2.Play(0, "joker", foe, "none")
	if len(tb2.G.Players[1].Hand) != 1 {
		t.Error("stealing is optional")
	}
}
