package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTumor(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("maggy")},
		Monsters: items("tumor", "fly"),
	})
	tb.G.DealDamageTo(engine.Target{Object: topID(tb, 0)}, 1, 0, 0)
	tb.Settle()
	if tb.G.Players[0].Damage != 0 || tb.G.Players[1].Damage != 1 || tb.G.Players[2].Damage != 1 {
		t.Error("each non-active player takes 1")
	}
}
