package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSoyMilk(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"soy_milk"}}, seat("cain")},
		Monsters: []engine.CardRef{"gurdy", "fatty"},
	})
	fatty, _ := tb.G.Monsters[1].TopOf()
	if dc := tb.G.Evasion(fatty); dc != 3 {
		t.Errorf("fatty DC = %d on Isaac's turn, want 2 + 1", dc)
	}
	tb.Attack("fatty", 3, 3, 3, 3) // each hit also deals 1 to Gurdy (HP 5)
	gurdy, _ := tb.G.Monsters[0].TopOf()
	if d := tb.G.Object(gurdy).Damage; d != 4 {
		t.Errorf("Gurdy damage = %d, want 4 (one per hit)", d)
	}
}

func TestSoyMilkOnlyYourTurn(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), {Character: "cain", Items: []engine.CardRef{"soy_milk"}}},
		Monsters: []engine.CardRef{"monstro", "fatty"},
	})
	fatty, _ := tb.G.Monsters[1].TopOf()
	if dc := tb.G.Evasion(fatty); dc != 2 {
		t.Errorf("fatty DC = %d on another player's turn, want 2", dc)
	}
}
