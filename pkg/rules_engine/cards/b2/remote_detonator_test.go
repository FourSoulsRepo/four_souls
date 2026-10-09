package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestRemoteDetonator(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "remote_detonator"), seat("cain", "breakfast"),
	}}, Set)
	tb.Activate(0, "remote_detonator", 0, "breakfast", "breakfast")
	if hasItem(tb.G, 1, "breakfast") {
		t.Error("the breakfast got two votes but stayed")
	}
	tb = enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "remote_detonator"), seat("cain", "breakfast"),
	}}, Set)
	tb.Activate(0, "remote_detonator", 0, "breakfast", "remote_detonator")
	if !hasItem(tb.G, 1, "breakfast") || !hasItem(tb.G, 0, "remote_detonator") {
		t.Error("a tie destroyed something")
	}
}
