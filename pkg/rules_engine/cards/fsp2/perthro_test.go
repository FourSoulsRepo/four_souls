package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPerthro(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"mystery_sack"}, Hand: []engine.CardRef{"perthro"}}, seat("cain"),
	}})
	tb.Play(0, "perthro", "mystery_sack")
	if hasItem(tb.G, 0, "mystery_sack") || len(tb.G.Players[0].InPlay) != 1 {
		t.Error("the sack should be rerolled into another treasure")
	}
}
