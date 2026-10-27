package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestEden(t *testing.T) {
	testExtraLoot(t, "eden")
	// Eden's choice comes before the first-player roll (R-SETUP-09).
	g, _, err := engine.NewGame(engine.Setup{Seed: 7, Players: 2, Sets: []engine.CardSet{Set}, Characters: []engine.CardRef{"isaac", "eden"}})
	if err != nil {
		t.Fatal(err)
	}
	p := g.Prompt()
	if p.Kind != engine.PromptChoose || p.Player != 1 || len(p.Options) != 3 {
		t.Fatalf("first prompt %+v, want Eden's player choosing among 3", p)
	}
	if g.FirstRoll == nil {
		t.Fatal("the first-player roll waits for Eden's choice")
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
	// No order question: the other two go to the bottom at random (Q-22).
	deck := g.Decks[engine.TreasureDeck]
	if len(deck) != bottom-1 {
		t.Errorf("treasure deck has %d cards, want %d", len(deck), bottom-1)
	}
	for i, id := range deck[:2] {
		if c := g.Object(id).Card; c != g.Object(top[0]).Card && c != g.Object(top[2]).Card {
			t.Errorf("bottom card %d is %s, not one of the other two", i, c)
		}
	}
	if g.Prompt().Kind != engine.PromptThrow || g.Prompt().Player != 0 {
		t.Errorf("after the choice the game waits on %+v, want the first roll", g.Prompt())
	}
}
