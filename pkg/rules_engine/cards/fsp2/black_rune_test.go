package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBlackRune(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"black_rune", "bomb"}}, seat("cain")},
		Monsters: []engine.CardRef{"monstro", "fatty"},
	})
	tb.G.ForceRolls(6)
	tb.Play(0, "black_rune")
	if n := len(tb.G.Players[0].Hand); n != 3 {
		t.Errorf("5-6: hand = %d, want 3", n)
	}
	tb2 := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"black_rune"}}, seat("cain")},
		Monsters: []engine.CardRef{"monstro", "fatty"},
	})
	tb2.G.ForceRolls(1)
	tb2.Play(0, "black_rune", "the rest in slot order")
	if d := tb2.G.Object(topID(tb2, 0)).Damage; d != 1 {
		t.Errorf("1-2: Monstro damage = %d, want 1", d)
	}
}
