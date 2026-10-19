package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAbyss(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("the_empty", "abyss", "spoon_bender", "jawbone"), seat("cain")}})
	tb.Activate(0, "abyss", 0, "spoon_bender")
	tb.Activate(0, "abyss", 0, "jawbone")
	if tb.G.PlayerATK(0) != 2 || len(tb.G.Players[0].InPlay) != 1 {
		t.Errorf("2 items destroyed: 2 counters, +1 ATK; ATK %d", tb.G.PlayerATK(0))
	}
}
