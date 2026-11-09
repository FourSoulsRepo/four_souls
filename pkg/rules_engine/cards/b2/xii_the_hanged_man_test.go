package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestXIITheHangedMan(t *testing.T) {
	tb := lootTable(t, "xii_the_hanged_man")
	treasure := tb.G.DeckTop(engine.TreasureDeck, 1)[0]
	loot := tb.G.DeckTop(engine.LootDeck, 1)[0]
	label := "ask.put_on_bottom " + string(tb.G.Object(treasure).Card)
	keep := func(id engine.ObjectID) string { return "ask.keep_on_top " + string(tb.G.Object(id).Card) }
	monster := tb.G.DeckTop(engine.MonsterDeck, 1)[0]
	tb.Play(0, "xii_the_hanged_man", label, keep(loot), keep(monster))
	if tb.G.Decks[engine.TreasureDeck][0] != treasure {
		t.Error("the treasure card is not on the bottom")
	}
	if n := len(tb.G.Players[0].Hand); n != 2 {
		t.Errorf("hand %d, want 2", n)
	}
}

// TestXIITheHangedManTellsEveryoneTheChoice: one top_choice per deck,
// in the deck order of its questions; each question shows others the
// look at its deck (19.147).
func TestXIITheHangedManTellsEveryoneTheChoice(t *testing.T) {
	tb := lootTable(t, "xii_the_hanged_man")
	decks := tb.G.DecksInPlay()
	tops := make([]engine.ObjectID, len(decks))
	for k, d := range decks {
		tops[k] = tb.G.DeckTop(d, 1)[0]
	}
	in := engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "xii_the_hanged_man")}}
	if _, err := tb.G.Apply(in); err != nil {
		t.Fatal(err)
	}
	tb.UntilChoose()
	var ev []engine.Event
	for k, d := range decks {
		look := engine.DeckLook{Deck: d, Count: 1}
		if l := tb.G.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 1}).Waiting.Look; l == nil || *l != look {
			t.Errorf("deck %s: others' look %+v, want %+v", d, l, look)
		}
		if s := tb.G.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 0}).Subject; s == nil || s.ID != tops[k] {
			t.Errorf("deck %s: chooser's subject %+v, want its top", d, s)
		}
		got, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentChoose, Choice: k % 2})
		if err != nil {
			t.Fatal(err)
		}
		ev = append(ev, got...)
	}
	var choices []engine.Event
	for _, e := range ev {
		if e.Kind == engine.EvTopChoice {
			choices = append(choices, e)
		}
	}
	if len(choices) != len(decks) {
		t.Fatalf("top_choice events %+v, want one per deck %v", choices, decks)
	}
	for k, d := range decks {
		where := " deck top"
		if k%2 == 1 {
			where = " deck bottom"
		}
		if e := choices[k]; e.Text != d.String()+where || e.Object != tops[k] || !e.Private {
			t.Errorf("deck %s: event %+v", d, e)
		}
	}
	if slices.ContainsFunc(ev, func(e engine.Event) bool { return e.Kind == engine.EvDeckBottom }) {
		t.Error("a deck_bottom too; top_choice tells it")
	}
}
