package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestTheChest(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "the_d20"), seat("cain", "the_chest")}}, Set)
	tb.Activate(0, "the_d20", 0, "the_chest")
	if tb.G.SoulValue(1) != 1 {
		t.Error("the chest did not become a soul")
	}
}
