package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestShadowOfWar(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac"), seat("cain", "spoon_bender", "jawbone")},
		Rooms:   []engine.CardRef{"shadow_of_war"},
	})
	tb.G.DealDamageTo(engine.Target{Player: 1, IsPlayer: true}, 9, 0, 0)
	tb.Settle("jawbone", "spoon_bender")
	if hasItem(tb.G, 1, "jawbone") || hasItem(tb.G, 1, "spoon_bender") {
		t.Error("Cain destroys an item before the penalty, then pays it")
	}
}
