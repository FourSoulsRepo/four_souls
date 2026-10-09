package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheCompass(t *testing.T) {
	tb := itemTable(t, items("the_compass"))
	var want []string
	for _, id := range tb.G.DeckTop(engine.LootDeck, 4) {
		want = append([]string{string(tb.G.Object(id).Card)}, want...)
	}
	tb.EndTurn(want...)
	var got []string
	for _, id := range tb.G.DeckTop(engine.LootDeck, 4) {
		got = append(got, string(tb.G.Object(id).Card))
	}
	// Cain's loot step then draws the new top card.
	if !slices.Equal(got[:3], want[1:]) {
		t.Errorf("top 3 = %v, want %v", got[:3], want[1:])
	}
}
