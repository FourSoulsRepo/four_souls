package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestPlacebo(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "placebo"), seat("cain", "mystery_sack", "the_d6")}}, Set)
	if n := len(tb.G.AbilitiesOf(tb.Find(0, "placebo"))); n != 1 {
		t.Fatalf("%d borrowed abilities, want 1 (the sack; the D6 is eternal)", n)
	}
	tb.G.ForceRolls(3)
	tb.Activate(0, "placebo", 0)
	if c := tb.G.Players[0].Cents; c != 4 || tb.G.Object(tb.Find(0, "placebo")).Charged {
		t.Errorf("cents %d: Placebo should use the sack's roll and deactivate", c)
	}
}
