package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAlabasterBox(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("isaac", "alabaster_box")}})
	box := tb.Find(1, "alabaster_box")
	tb.G.AddCountersTo(box, 1)
	tb.EndTurn() // Isaac's turn: the 2nd counter
	if got := tb.G.Players[1].Cents; got != 10 {
		t.Errorf("cents = %d, want 10 at the 2nd counter", got)
	}
	tb.G.AddCountersTo(box, 1)
	tb.EndTurn()
	tb.EndTurn() // the 4th counter
	if got := len(tb.G.Players[1].Hand); got < 4 {
		t.Errorf("hand = %d, want 4 more cards at the 4th counter", got)
	}
}
