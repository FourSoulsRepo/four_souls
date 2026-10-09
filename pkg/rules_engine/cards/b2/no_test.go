package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestNo(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "jawbone"), {Character: "cain", Items: items("no"), Cents: 3},
	}}, Set)
	tb.Start(0, "jawbone", 0)
	tb.Choose(foe) // the steal waits on the stack
	tb.Pass(0)
	tb.Activate(1, "no", 0, "jawbone: ↷: Steal 3¢ from a player.")
	if c := tb.G.Players[1].Cents; c != 3 {
		t.Errorf("Cain has %d¢, want 3: the steal was cancelled", c)
	}
}
