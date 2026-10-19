package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPageantBoy(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "pageant_boy"), seat("cain")}})
	tb.G.ForceRolls(5)
	tb.Activate(0, "pageant_boy", 0)
	if hasItem(tb.G, 0, "pageant_boy") || tb.G.Players[0].Cents != 20 {
		t.Errorf("destroyed for 4 × 5 = 20¢; cents %d", tb.G.Players[0].Cents)
	}
}
