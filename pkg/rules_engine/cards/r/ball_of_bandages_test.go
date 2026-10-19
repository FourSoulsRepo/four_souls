package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBallOfBandages(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "ball_of_bandages"), seat("cain")}, Monsters: []engine.CardRef{"pooter"}})
	tb.Attack("pooter", 3, 6, 6)
	if tb.G.Players[0].Damage != 0 {
		t.Error("no combat damage on an attack roll of 3")
	}
}
