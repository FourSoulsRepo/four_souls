package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMomsEyeShadow(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), {Character: "cain", Items: items("moms_eye_shadow")}, seat("maggy")},
		Monsters: items("fly", "gurdy"),
	})
	tb.G.ForceRolls(1, 6, 6, 6, 6, 6)
	tb.Do(attackIntent(tb), "fly", "gurdy") // Cain sends Isaac against Gurdy
	if tb.G.Players[0].Damage != 1 {
		t.Errorf("Isaac's damage = %d: he missed Gurdy (ATK 1)", tb.G.Players[0].Damage)
	}
	if top(tb, 0) != "fly" {
		t.Error("the fly was not attacked")
	}
}
