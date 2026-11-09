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
	// The roll starts only after her choice: no roll in the view yet (19.142).
	if g.FirstRoll != nil {
		t.Fatalf("first-player roll %+v during Eden's choice, want none", g.FirstRoll)
	}
	if v := g.View(engine.Viewer{Player: 0}); v.FirstRoll != nil {
		t.Fatalf("the view sends the roll %+v during Eden's choice", v.FirstRoll)
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
	if g.Prompt().Kind != engine.PromptThrow || g.Prompt().Player != 0 || g.FirstRoll == nil || g.FirstRoll.Round != 1 {
		t.Errorf("after the choice the game waits on %+v, roll %+v, want the first roll", g.Prompt(), g.FirstRoll)
	}
}

// TestEdenBeforeCain: with Cain at the table, Eden still chooses first;
// Cain's start comes after her choice (R-SETUP-09, 19.142).
func TestEdenBeforeCain(t *testing.T) {
	g, ev, err := engine.NewGame(engine.Setup{Seed: 3, Players: 2, Sets: []engine.CardSet{Set}, Characters: []engine.CardRef{"cain", "eden"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range ev {
		if e.Kind == engine.EvFirstPlayer || e.Kind == engine.EvGameStarted {
			t.Fatalf("event %s before Eden's choice", e.Kind)
		}
	}
	if p := g.Prompt(); p.Kind != engine.PromptChoose || p.Player != 1 || g.Turn.Number != 0 {
		t.Fatalf("prompt %+v in turn %d, want Eden's choice before turn 1", p, g.Turn.Number)
	}
	ev, err = g.Apply(engine.Intent{Player: 1, Kind: engine.IntentChoose, Choice: 0})
	if err != nil {
		t.Fatal(err)
	}
	started := false
	for _, e := range ev {
		started = started || (e.Kind == engine.EvGameStarted && e.Player == 0)
	}
	if !started || g.Turn.Number != 1 || g.Turn.Active != 0 {
		t.Fatalf("after the choice: turn %+v, started %v, want Cain's turn 1", g.Turn, started)
	}
}
