package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPotatoPeeler(t *testing.T) {
	tb := itemTable(t, items("potato_peeler"))
	tb.Activate(0, "potato_peeler", 0)
	for _, d := range []engine.DeckKind{engine.TreasureDeck, engine.LootDeck, engine.MonsterDeck} {
		if len(tb.G.Discards[d]) != 1 {
			t.Errorf("%s discard has %d cards, want 1", d, len(tb.G.Discards[d]))
		}
	}
}
