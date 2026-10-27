package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSackHead(t *testing.T) {
	tb := itemTable(t, items("sack_head"))
	top := tb.G.DeckTop(engine.LootDeck, 1)[0]
	tb.Activate(0, "sack_head", 0, "ask.deck.loot", "ask.put_on_bottom "+string(tb.G.Object(top).Card))
	if tb.G.Decks[engine.LootDeck][0] != top {
		t.Error("the card is not on the bottom")
	}
}
