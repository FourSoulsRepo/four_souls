package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestEden(t *testing.T) {
	testExtraLoot(t, "eden")
	g := newGame(t, "isaac", "eden")
	p := g.Prompt()
	if p.Kind != engine.PromptChoose || p.Player != 1 || len(p.Options) != 3 {
		t.Fatalf("first prompt %+v, want Eden's player choosing among 3", p)
	}
	top := g.DeckTop(engine.TreasureDeck, 3)
	bottom := len(g.Decks[engine.TreasureDeck])
	if _, err := g.Apply(engine.Intent{Player: 1, Kind: engine.IntentChoose, Choice: 1}); err != nil {
		t.Fatal(err)
	}
	chosen := g.Object(top[1]).Card
	var item *engine.Object
	for _, id := range g.Players[1].InPlay {
		if g.Object(id).Card == chosen {
			item = g.Object(id)
		}
	}
	if item == nil || !item.Eternal || !item.Charged {
		t.Fatalf("starting item %s: %+v, want an eternal charged item", chosen, item)
	}
	deck := g.Decks[engine.TreasureDeck]
	if len(deck) != bottom-1 {
		t.Errorf("treasure deck has %d cards, want %d", len(deck), bottom-1)
	}
	for i, id := range deck[:2] {
		if c := g.Object(id).Card; c != g.Object(top[0]).Card && c != g.Object(top[2]).Card {
			t.Errorf("bottom card %d is %s, not one of the other two", i, c)
		}
	}
	if g.Prompt().Kind != engine.PromptPriority {
		t.Errorf("after the choice the game waits on %+v", g.Prompt())
	}
}
