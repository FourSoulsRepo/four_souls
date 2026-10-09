package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestRestock(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "restock"), seat("cain")},
		Shop:    items("breakfast", "dinner"),
		Active:  1,
	}, Set)
	tb.EndTurn("discard breakfast", "keep dinner")
	if !slices.ContainsFunc(tb.G.Discards[engine.TreasureDeck], func(id engine.ObjectID) bool { return tb.G.Object(id).Card == "breakfast" }) {
		t.Error("breakfast was not discarded")
	}
}
