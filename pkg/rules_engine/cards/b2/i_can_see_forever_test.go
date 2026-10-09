package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestICanSeeForever(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	var order []string
	for _, id := range tb.G.DeckTop(engine.LootDeck, 6) {
		order = append([]string{string(tb.G.Object(id).Card)}, order...)
	}
	tb.RevealFromDeck("i_can_see_forever", order...)
	if h := handCards(tb.G, 0); len(h) != 1 || string(h[0]) != order[0] {
		t.Errorf("hand %v, want the new top card %s", h, order[0])
	}
}
