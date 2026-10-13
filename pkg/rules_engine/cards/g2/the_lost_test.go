package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheLost(t *testing.T) {
	testCharacter(t, "the_lost", "holy_mantle")
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "the_lost"},
		{Character: "cain", Deactivated: []engine.CardRef{"character"}},
	}})
	if hp := tb.G.PlayerHP(0); hp != 1 {
		t.Errorf("HP = %d, want 1", hp)
	}
	if v := tb.G.SoulValue(0); v != 1 {
		t.Errorf("soul value = %d, want 1 (the character is a soul)", v)
	}
	tb.EndTurn("cain")
	if !tb.G.Object(tb.G.Players[1].Character).Charged {
		t.Error("Cain was not recharged at the end of the Lost's turn")
	}
}
