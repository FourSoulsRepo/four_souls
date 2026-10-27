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
	tb.Activate(0, "sleight_of_hand", 0, append([]string{"ask.deck.monster"}, reversed...)...)
	var got []string
	for _, id := range tb.G.DeckTop(engine.MonsterDeck, 5) {
		got = append(got, string(tb.G.Object(id).Card))
	}
	if !slices.Equal(got, reversed) {
		t.Errorf("top 5 = %v, want %v", got, reversed)
	}
}

// TestSleightOfHandAsksOneOrder: the five cards come in one order prompt
// that only the chooser sees; the answer must hold each card once (19.16).
func TestSleightOfHandAsksOneOrder(t *testing.T) {
	tb := table(t, seat("cain", "sleight_of_hand"), seat("isaac"))
	top := tb.G.DeckTop(engine.MonsterDeck, 5)
	tb.Start(0, "sleight_of_hand", 0)
	for tb.G.Prompt().Kind == engine.PromptPriority {
		tb.Pass(int(tb.G.Prompt().Player))
	}
	tb.Choose("ask.deck.monster")
	w := tb.G.Prompt()
	if w.Kind != engine.PromptOrder || w.Count != 5 || w.Text.String() != "ask.put_back" {
		t.Fatalf("prompt = %+v, want an order of 5", w)
	}
	var shown []engine.ObjectID
	for _, c := range tb.G.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 0}).Ordering {
		shown = append(shown, c.ID)
	}
	if !slices.Equal(shown, top) {
		t.Errorf("chooser sees %v, want the top 5 %v", shown, top)
	}
	for _, v := range []engine.Viewer{{Kind: engine.ViewPlayer, Player: 1}, {Kind: engine.ViewSpectator}} {
		if o := tb.G.View(v).Ordering; o != nil {
			t.Errorf("viewer %+v sees the cards %v", v, o)
		}
	}
	allowed := tb.G.Allowed(0)
	if len(allowed) != 1 || allowed[0].Kind != engine.IntentOrder || !slices.Equal(allowed[0].Objects, top) {
		t.Fatalf("allowed = %+v, want one order intent with the top 5", allowed)
	}
	for _, bad := range [][]engine.ObjectID{top[:4], {top[0], top[0], top[1], top[2], top[3]}, append(slices.Clone(top[:4]), tb.Find(0, "cain"))} {
		if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentOrder, Objects: bad}); err == nil {
			t.Errorf("order %v accepted", bad)
		}
	}
	want := []engine.ObjectID{top[2], top[0], top[4], top[1], top[3]}
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentOrder, Objects: want})
	if got := tb.G.DeckTop(engine.MonsterDeck, 5); !slices.Equal(got, want) {
		t.Errorf("top 5 = %v, want %v", got, want)
	}
}
