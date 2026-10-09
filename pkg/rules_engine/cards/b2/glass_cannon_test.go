package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestGlassCannon(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "glass_cannon"), seat("cain", "breakfast")}}, Set)
	tb.G.ForceRolls(6)
	tb.Activate(0, "glass_cannon", 0, "breakfast")
	if hasItem(tb.G, 1, "breakfast") || !tb.G.Object(tb.Find(0, "glass_cannon")).Charged {
		t.Error("breakfast should be gone and the cannon recharged")
	}
	tb = enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "glass_cannon"), seat("cain", "breakfast")}}, Set)
	tb.G.ForceRolls(1)
	tb.Activate(0, "glass_cannon", 0, "breakfast")
	if hasItem(tb.G, 0, "glass_cannon") || len(tb.G.Players[0].Hand) != 2 {
		t.Error("on 1-5 the cannon breaks and you loot 2")
	}
}
