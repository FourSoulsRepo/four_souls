package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestTheCurse(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("eve", "the_curse"), seat("cain")},
		Active:  1,
	}, Set)
	lootDeck := len(tb.G.Decks[engine.LootDeck])
	tb.EndTurn("loot deck") // Eve's turn starts: the trigger asks for a deck
	if got := len(tb.G.Discards[engine.LootDeck]); got != 1 {
		t.Fatalf("loot discard has %d cards, want 1 (milled)", got)
	}
	if got := len(tb.G.Decks[engine.LootDeck]); got != lootDeck-2 {
		t.Errorf("loot deck has %d cards, want %d (milled one, looted one)", got, lootDeck-2)
	}

	milled := tb.G.Object(tb.G.Discards[engine.LootDeck][0]).Card
	tb.Activate(0, "the_curse", 1, "loot discard")
	if len(tb.G.Discards[engine.LootDeck]) != 0 {
		t.Error("the discard still has the card")
	}
	if top := tb.G.DeckTop(engine.LootDeck, 1); len(top) != 1 || tb.G.Object(top[0]).Card != milled {
		t.Errorf("top of the loot deck is not the milled %s", milled)
	}
}
