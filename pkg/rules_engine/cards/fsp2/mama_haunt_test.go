package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMamaHaunt(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac"), {Character: "cain", Items: items("mama_haunt"), Deactivated: items("character")},
	}})
	tb.EndTurn()
	if tb.G.Object(tb.G.Players[1].Character).Charged {
		t.Error("Cain's character doesn't recharge")
	}
	tb.G.DealDamageTo(engine.Target{Player: 1, IsPlayer: true}, 2, 0, 0)
	settleAny(tb)
	if !hasItem(tb.G, 0, "mama_haunt") {
		t.Error("on death it goes to another player")
	}
}
