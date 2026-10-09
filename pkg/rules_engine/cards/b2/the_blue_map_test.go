package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheBlueMap(t *testing.T) {
	tb := itemTable(t, items("the_blue_map"))
	var want []string
	for _, id := range tb.G.DeckTop(engine.TreasureDeck, 4) {
		want = append([]string{string(tb.G.Object(id).Card)}, want...)
	}
	tb.EndTurn(want...)
	var got []string
	for _, id := range tb.G.DeckTop(engine.TreasureDeck, 4) {
		got = append(got, string(tb.G.Object(id).Card))
	}
	if !slices.Equal(got, want) {
		t.Errorf("top 4 = %v, want %v", got, want)
	}
}
