package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGuppysTail(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"guppys_tail"}, Hand: []engine.CardRef{"bomb", "a_sack"}},
		seat("cain"),
	}})
	tb.EndTurn("bomb")
	if n := len(tb.G.Players[0].Hand); n != 1 {
		t.Errorf("hand = %d after the end of turn, want 1", n)
	}
	// Cain's turn, then Isaac's loot step gives 2 cards.
	tb.EndTurn("a_sack")
	if n := len(tb.G.Players[0].Hand); n != 3 {
		t.Errorf("hand = %d after the loot step, want 1 + 2", n)
	}
}
