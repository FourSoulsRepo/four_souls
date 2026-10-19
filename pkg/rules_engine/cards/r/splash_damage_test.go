package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSplashDamage(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"leech", "fly"},
		Rooms:    []engine.CardRef{"splash_damage"},
	})
	tb.Attack("leech", 6)
	if tb.G.Players[1].Damage != 1 || tb.G.Players[2].Damage != 0 {
		t.Errorf("Isaac's hit splashes 1 damage onto Cain, to his left: %d, %d", tb.G.Players[1].Damage, tb.G.Players[2].Damage)
	}
}
