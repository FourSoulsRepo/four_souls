package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestMomsRazor(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "moms_razor"), seat("cain", "mystery_sack")}, Active: 1}, Set)
	tb.G.ForceRolls(6)
	tb.Activate(1, "mystery_sack", 0, "yes")
	if d := tb.G.Players[1].Damage; d != 1 {
		t.Errorf("damage = %d, want 1", d)
	}
}
