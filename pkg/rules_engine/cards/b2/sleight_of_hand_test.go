package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSleightOfHand(t *testing.T) {
	tb := table(t, seat("cain", "sleight_of_hand"), seat("isaac"))
	var names []string
	for _, id := range tb.G.DeckTop(engine.MonsterDeck, 5) {
		names = append(names, string(tb.G.Object(id).Card))
	}
	reversed := slices.Clone(names)
	slices.Reverse(reversed)
	tb.Activate(0, "sleight_of_hand", 0, append([]string{"monster deck"}, reversed...)...)
	var got []string
	for _, id := range tb.G.DeckTop(engine.MonsterDeck, 5) {
		got = append(got, string(tb.G.Object(id).Card))
	}
	if !slices.Equal(got, reversed) {
		t.Errorf("top 5 = %v, want %v", got, reversed)
	}
}
