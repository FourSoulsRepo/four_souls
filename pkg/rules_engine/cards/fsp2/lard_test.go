package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLard(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("lard"), Hand: items("bomb")}, seat("cain"),
	}})
	if hp := tb.G.PlayerHP(0); hp != 4 {
		t.Errorf("HP = %d, want 2 + 2", hp)
	}
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 1, 1, 0)
	tb.Settle("bomb")
	if len(tb.G.Players[0].Hand) != 0 {
		t.Error("damage makes the player discard")
	}
}
