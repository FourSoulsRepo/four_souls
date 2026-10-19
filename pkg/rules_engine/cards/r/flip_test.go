package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestFlip(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"mega_bomb"}},
		{Character: "the_enigma", Items: []engine.CardRef{"flip"}, Hand: []engine.CardRef{"bomb", "bomb"}},
	}})
	tb.Play(0, "mega_bomb", "player 2 (the_enigma)", "bomb")
	ch := tb.G.Players[1].Character
	if tb.G.CardOf(ch) != "amgine_eht" || !tb.G.Object(ch).Charged {
		t.Error("after the penalty his character flips and recharges")
	}
	if len(tb.G.Players[1].Hand) != 3 {
		t.Errorf("hand %d: he discards it and loots 3", len(tb.G.Players[1].Hand))
	}
}
