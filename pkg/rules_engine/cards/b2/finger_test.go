package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestFinger(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "finger", "breakfast"), seat("cain", "mystery_sack", "dinner"),
	}, Active: 1}, Set)
	tb.G.ForceRolls(2)
	tb.Activate(1, "mystery_sack", 0, "breakfast", "dinner")
	if !hasItem(tb.G, 0, "dinner") || !hasItem(tb.G, 1, "breakfast") {
		t.Error("not swapped")
	}
}
