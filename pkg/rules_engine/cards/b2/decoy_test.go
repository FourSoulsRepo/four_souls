package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestDecoy(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "decoy"), seat("cain", "breakfast", "sleight_of_hand"),
	}}, Set)
	tb.Start(0, "decoy", 0)
	if opts := tb.Options(); len(opts) != 2 {
		t.Fatalf("options %v, want the breakfast and cancel (not the eternal item)", opts)
	}
	tb.Settle("breakfast")
	if !hasItem(tb.G, 0, "breakfast") || !hasItem(tb.G, 1, "decoy") {
		t.Error("not swapped")
	}
}
