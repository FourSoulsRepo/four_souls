package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestTheD20(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "the_d20"), seat("cain", "breakfast")}}, Set)
	tb.Activate(0, "the_d20", 0, "breakfast")
	if hasItem(tb.G, 1, "breakfast") || len(tb.G.Players[1].InPlay) != 1 {
		t.Error("the breakfast was not replaced by a new treasure")
	}
}
