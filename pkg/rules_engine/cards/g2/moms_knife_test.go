package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMomsKnife(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"moms_knife"}}, seat("cain")},
		Monsters: []engine.CardRef{"monstro", "fatty"},
	})
	if atk := tb.G.PlayerATK(0); atk != 2 {
		t.Errorf("ATK = %d, want 1 + 1", atk)
	}
	tb.Attack("fatty", 6, 2) // a 6 hits: 2 damage to the fatty, 1 to Monstro
	monstro, _ := tb.G.Monsters[0].TopOf()
	if d := tb.G.Object(monstro).Damage; d != 1 {
		t.Errorf("Monstro damage = %d, want 1 (only the roll of 6)", d)
	}
}
