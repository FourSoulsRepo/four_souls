package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMonstroII(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("maggy")},
		Monsters: items("monstro_ii", "fly"),
	})
	tb.G.ForceRolls(1, 5, 6, 6, 6, 6, 6) // Monstro II has 5 HP // a miss; Monstro II rolls 5: the player to the left (Cain) takes 1
	tb.Do(attackIntent(tb), "monstro_ii")
	if d := tb.G.Players[1].Damage; d != 1 {
		t.Errorf("Cain's damage = %d, want 1 (left = next in turn order)", d)
	}
}
