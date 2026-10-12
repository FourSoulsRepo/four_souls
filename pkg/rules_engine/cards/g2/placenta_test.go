package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPlacenta(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"placenta"}}, seat("cain")},
		Monsters: []engine.CardRef{"trite", "fly"},
	})
	if hp := tb.G.PlayerHP(0); hp != 3 {
		t.Errorf("HP = %d, want 2 + 1", hp)
	}
	tb.Attack("trite", 4, 5) // a 4 misses the trite (DC 5): no damage
	if d := tb.G.Players[0].Damage; d != 0 {
		t.Errorf("damage = %d, want 0 on an attack roll of 4", d)
	}
}
