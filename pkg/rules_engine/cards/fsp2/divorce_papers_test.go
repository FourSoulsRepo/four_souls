package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDivorcePapers(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("divorce_papers")},
		{Character: "cain", Cents: 7, Hand: items("bomb", "two_cents", "four_cents"), Items: items("mystery_sack")},
	}})
	tb.Activate(0, "divorce_papers", 0, foe, "bomb", "mystery_sack")
	if c := tb.G.Players[0].Cents; c != 3 {
		t.Errorf("cents = %d, want 7/2 = 3", c)
	}
	if n := len(tb.G.Players[0].Hand); n != 1 {
		t.Errorf("hand = %d, want 3/2 = 1", n)
	}
	if !hasItem(tb.G, 0, "mystery_sack") {
		t.Error("Cain gives an item")
	}
}
