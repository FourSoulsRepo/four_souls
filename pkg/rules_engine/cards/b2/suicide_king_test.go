package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSuicideKing(t *testing.T) {
	tb := itemTable(t, items("suicide_king"), "xiii_death")
	first := string(tb.G.Object(tb.G.DeckTop(engine.LootDeck, 1)[0]).Card)
	tb.Play(0, "xiii_death", me, "suicide_king", first)
	// Loot 3 before the penalty, then discard one of them.
	if h := len(tb.G.Players[0].Hand); h != 2 {
		t.Errorf("hand = %d, want 2", h)
	}
}
