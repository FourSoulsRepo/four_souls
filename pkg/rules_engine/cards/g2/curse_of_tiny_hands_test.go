package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCurseOfTinyHands(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"bomb", "bomb", "a_sack", "jera"}}, seat("cain")},
		Monsters: []engine.CardRef{"fatty", "monstro"},
	})
	tb.RevealFromDeck("curse_of_tiny_hands", me)
	tb.EndTurn()
	if n := len(tb.G.Players[0].Hand); n != 2 {
		t.Errorf("hand = %d at the end of the turn, want 2", n)
	}
}
