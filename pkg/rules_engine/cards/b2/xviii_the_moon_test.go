package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestXVIIITheMoon(t *testing.T) { testOneOnTop(t, "xviii_the_moon", engine.LootDeck) }

// TestXVIIITheMoonShowsTheLook: while the chooser picks, every viewer
// sees which deck is looked at and how many cards; only the chooser sees
// the cards. The look ends with the answer (20.18).
func TestXVIIITheMoonShowsTheLook(t *testing.T) {
	tb := lootTable(t, "xviii_the_moon")
	in := engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "xviii_the_moon")}}
	if _, err := tb.G.Apply(in); err != nil {
		t.Fatal(err)
	}
	for tb.G.Prompt().Kind == engine.PromptPriority {
		tb.Pass(int(tb.G.Prompt().Player))
	}
	want := engine.DeckLook{Deck: engine.LootDeck, Count: 5}
	if w := tb.G.Prompt(); w.Kind != engine.PromptChoose || w.Look == nil || *w.Look != want {
		t.Fatalf("prompt = %+v, want a choice with the look %+v", w, want)
	}
	for _, v := range []engine.Viewer{{Kind: engine.ViewPlayer, Player: 0}, {Kind: engine.ViewPlayer, Player: 1}, {Kind: engine.ViewSpectator}} {
		view := tb.G.View(v)
		if l := view.Waiting.Look; l == nil || *l != want {
			t.Errorf("viewer %+v: look %+v, want %+v", v, l, want)
		}
		if v.Player == 1 && view.Waiting.Options != nil {
			t.Errorf("viewer %+v sees the options %v", v, view.Waiting.Options)
		}
	}
	tb.Choose(tb.Options()[0])
	if l := tb.G.Prompt().Look; l != nil {
		t.Errorf("look %+v after the answer, want none", l)
	}
}

// TestXVIIITheMoonTellsEveryoneTheBottom: the 4 cards not kept go to
// the bottom as one public deck_bottom, with no card (19.147).
func TestXVIIITheMoonTellsEveryoneTheBottom(t *testing.T) {
	tb := lootTable(t, "xviii_the_moon")
	ev := tb.Play(0, "xviii_the_moon", string(tb.G.Object(tb.G.DeckTop(engine.LootDeck, 1)[0]).Card))
	var got []engine.Event
	for _, e := range ev {
		if e.Kind == engine.EvDeckBottom {
			got = append(got, e)
		}
	}
	if len(got) != 1 {
		t.Fatalf("deck_bottom events %+v, want one", got)
	}
	if e := got[0]; e.Amount != 4 || e.Text != "loot deck" || e.Private || e.Card != "" || e.Player != engine.NoPlayer {
		t.Errorf("event %+v", e)
	}
	if engine.Barrier(got[0]) {
		t.Error("deck_bottom is a barrier")
	}
}
