package b2

import (
	"slices"
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

// TestSackHeadTellsEveryoneTheChoice: Sack Head's answer is a
// top_choice, as Cain's Eye's; others see the look while it waits
// (19.147).
func TestSackHeadTellsEveryoneTheChoice(t *testing.T) {
	for choice, want := range []string{"loot deck top", "loot deck bottom"} {
		tb := itemTable(t, items("sack_head"))
		top := tb.G.DeckTop(engine.LootDeck, 1)[0]
		tb.Start(0, "sack_head", 0)
		tb.UntilChoose()
		tb.Choose("ask.deck.loot")
		look := engine.DeckLook{Deck: engine.LootDeck, Count: 1}
		if l := tb.G.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 1}).Waiting.Look; l == nil || *l != look {
			t.Errorf("%s: others' look %+v, want %+v", want, l, look)
		}
		ev, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentChoose, Choice: choice})
		if err != nil {
			t.Fatal(err)
		}
		i := slices.IndexFunc(ev, func(e engine.Event) bool { return e.Kind == engine.EvTopChoice })
		if i < 0 {
			t.Fatalf("%s: no top_choice in %+v", want, ev)
		}
		if e := ev[i]; e.Player != 0 || e.Text != want || e.Object != top || !e.Private {
			t.Errorf("%s: event %+v", want, e)
		}
	}
}
