package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestMonstrosTooth(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "monstros_tooth", "breakfast"), seat("cain", "breakfast")}, Active: 1}, Set)
	before := len(tb.G.Players[0].InPlay) + len(tb.G.Players[1].InPlay)
	// Either player may be picked; both have a Breakfast.
	tb.EndTurn("breakfast")
	if after := len(tb.G.Players[0].InPlay) + len(tb.G.Players[1].InPlay); after != before-1 {
		t.Errorf("items %d, want %d", after, before-1)
	}
}
