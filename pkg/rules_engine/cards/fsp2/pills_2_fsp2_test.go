package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPills2Fsp2(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"mystery_sack", "the_d6"}, Hand: []engine.CardRef{"pills_2_fsp2"}}, seat("cain"),
	}})
	tb.G.ForceRolls(5)
	tb.Play(0, "pills_2_fsp2")
	if hasItem(tb.G, 0, "mystery_sack") || !hasItem(tb.G, 0, "the_d6") {
		t.Error("5-6: each non-eternal item is rerolled; the eternal D6 stays")
	}
}
