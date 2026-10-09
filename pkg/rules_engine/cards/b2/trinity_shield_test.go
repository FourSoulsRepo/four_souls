package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestTrinityShield(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "trinity_shield", "jawbone"), seat("cain", "mystery_sack"),
	}}, Set)
	tb.Start(0, "jawbone", 0)
	tb.Choose(foe)
	tb.Pass(0)
	if rule := tb.Refused(1, "mystery_sack", 0); rule != "R-ABIL-12" {
		t.Errorf("refused by %s, want R-ABIL-12", rule)
	}
	tb.Activate(1, "cain", 0) // characters are not items
}
