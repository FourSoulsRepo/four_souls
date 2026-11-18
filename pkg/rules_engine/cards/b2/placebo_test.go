package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// TestPlaceboBorrowsNoEternalAbility: Placebo copies only a non-eternal
// item's ↷ (an engine API check). Its use is a situation file:
// placebo/uses-a-copied-roll-ability.json.
func TestPlaceboBorrowsNoEternalAbility(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "placebo"), seat("cain", "mystery_sack", "the_d6")}}, Set)
	if n := len(tb.G.AbilitiesOf(tb.Find(0, "placebo"))); n != 1 {
		t.Fatalf("%d borrowed abilities, want 1 (the sack; the D6 is eternal)", n)
	}
}
