package r

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestAnimaSolaShowsTheRevenant: the view shows the face up (ADR 023).
// The flip itself is in testdata/situations/anima_sola/.
func TestAnimaSolaShowsTheRevenant(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("the_deserter", "anima_sola"), seat("cain")}})
	sola := tb.Find(0, "anima_sola")
	tb.Activate(0, "anima_sola", 0)
	in := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Players[0].InPlay
	if i := slices.IndexFunc(in, func(cv engine.CardView) bool { return cv.ID == sola }); i < 0 || in[i].Card != "the_revenant" || !in[i].Flipped {
		t.Errorf("in play %+v: the face up is The Revenant (ADR 023)", in)
	}
}
