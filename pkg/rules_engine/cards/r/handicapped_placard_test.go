package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHandicappedPlacard(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "handicapped_placard"), {Character: "cain", Souls: []engine.CardRef{"monstro", "gurdy"}}},
		Monsters: []engine.CardRef{"leech"},
	})
	if dc := tb.G.Evasion(tb.G.Monsters[0].Cards[0]); dc != 2 {
		t.Errorf("DC = %d, want 4 - 2 on his turn", dc)
	}
}
