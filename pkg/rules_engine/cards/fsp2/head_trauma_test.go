package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHeadTrauma(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: items("bomb", "two_cents", "four_cents")}, seat("cain")},
		Monsters: items("fatty", "monstro"),
	})
	tb.RevealFromDeck("head_trauma", "bomb", "two_cents")
	if n := len(tb.G.Players[0].Hand); n != 1 {
		t.Errorf("hand = %d, want 1", n)
	}
}
