package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestTheD4(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "the_d4"), seat("cain", "breakfast", "dinner")}}, Set)
	tb.Activate(0, "the_d4", 0, foe)
	if hasItem(tb.G, 1, "breakfast") || hasItem(tb.G, 1, "dinner") || len(tb.G.Players[1].InPlay) != 2 {
		t.Error("Cain's items were not rerolled")
	}
}
