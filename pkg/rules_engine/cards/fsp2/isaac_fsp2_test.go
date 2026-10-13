package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestIsaacFsp2(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("maggy")},
		Monsters: items("isaac_fsp2", "fly"),
	})
	tb.G.DealDamageTo(engine.Target{Object: topID(tb, 0)}, 1, 0, 0)
	tb.Settle("player 3 (maggy)")
	if d := tb.G.Players[2].Damage; d != 1 {
		t.Errorf("Maggy's damage = %d, want 1", d)
	}
}
