package rulesengine

import (
	"slices"
	"testing"
)

// deckGame is a bare game with cards a to e in the loot deck, e on top.
func deckGame() *Game {
	g := &Game{RNG: NewRNG(1)}
	g.AddToDeck(LootDeck, "a", "b", "c", "d", "e")
	return g
}

// TestLookPlacesDeckCard: a look at deck cards names each card's deck
// and depth; FilterEvents hides them with the card from other players
// and shows them to the looker and a judge (ADR 030, Decision 19.11).
func TestLookPlacesDeckCard(t *testing.T) {
	g := deckGame()
	g.LookAt(0, g.DeckTop(LootDeck, 2)...)
	events := g.takeEvents()
	if len(events) != 2 {
		t.Fatalf("events %+v, want two looks", events)
	}
	for i, e := range events {
		if e.Deck == nil || e.Deck.Deck != LootDeck || e.Deck.Depth != i {
			t.Errorf("look %d: place %+v, want loot deck depth %d", i, e.Deck, i)
		}
	}
	for _, v := range []Viewer{{Kind: ViewPlayer, Player: 0}, {Kind: ViewJudge}} {
		for _, e := range FilterEvents(events, v) {
			if e.Deck == nil || e.Card == "" {
				t.Errorf("viewer %+v: %+v lost its card or place", v, e)
			}
		}
	}
	for _, v := range []Viewer{{Kind: ViewPlayer, Player: 1}, {Kind: ViewSpectator}} {
		for _, e := range FilterEvents(events, v) {
			if e.Deck != nil || e.Card != "" || e.Object != 0 {
				t.Errorf("viewer %+v sees %+v", v, e)
			}
		}
	}
}

// TestDeckOrderEvents: a shuffle, a reorder of the top cards, a card
// moved to the bottom and a refill from the discard tell deck_order; a
// kept or bottomed top card and cards put under the deck tell only
// their own event, and setup shuffles tell nothing.
func TestDeckOrderEvents(t *testing.T) {
	orders := func(events []Event) int {
		n := 0
		for _, e := range events {
			if e.Kind == EvDeckOrder {
				if e.Text != "loot deck" || e.Private || e.Card != "" {
					t.Errorf("deck_order %+v", e)
				}
				n++
			}
		}
		return n
	}
	for name, tc := range map[string]struct {
		do   func(g *Game)
		want int
	}{
		"shuffle":       {func(g *Game) { g.ShuffleDeck(LootDeck) }, 1},
		"setup shuffle": {func(g *Game) { g.shuffleDeck(LootDeck) }, 0},
		"reorder": {func(g *Game) {
			top := g.DeckTop(LootDeck, 2)
			g.SetDeckTop(LootDeck, []ObjectID{top[1], top[0]})
		}, 1},
		"to the bottom":  {func(g *Game) { g.DeckToBottom(LootDeck, g.DeckTop(LootDeck, 1)[0]) }, 1},
		"keep or bottom": {func(g *Game) { g.KeepOrBottom(0, LootDeck, g.DeckTop(LootDeck, 1)[0], true) }, 0},
		"bottom random":  {func(g *Game) { g.ToBottomRandom(LootDeck, g.DeckTop(LootDeck, 2)) }, 0},
		"refill": {func(g *Game) {
			for range 5 {
				g.MillTop(LootDeck)
			}
			g.drawTop(LootDeck)
		}, 1},
	} {
		t.Run(name, func(t *testing.T) {
			g := deckGame()
			before := slices.Clone(g.Decks[LootDeck])
			tc.do(g)
			if got := orders(g.takeEvents()); got != tc.want {
				t.Errorf("%d deck_order events, want %d", got, tc.want)
			}
			if name == "setup shuffle" && len(g.Decks[LootDeck]) != len(before) {
				t.Error("the shuffle lost cards")
			}
		})
	}
}

// TestDeckText: every deck's event text reads back as that deck, with
// its place; a text with another place, or none of a deck, does not.
func TestDeckText(t *testing.T) {
	wheres := []string{"", DeckTextTop, DeckTextBottom}
	for d := range deckCount {
		for _, w := range wheres {
			text := DeckText(d, w)
			if got, ok := DeckOfText(text, w); !ok || got != d {
				t.Errorf("DeckOfText(%q, %q) = %v, %v; want %v", text, w, got, ok, d)
			}
			for _, other := range wheres {
				if other == w {
					continue
				}
				if _, ok := DeckOfText(text, other); ok {
					t.Errorf("DeckOfText(%q, %q) reads a deck", text, other)
				}
			}
		}
	}
	for _, text := range []string{"", "deck", TextFromDiscard, "shop deck"} {
		if _, ok := DeckOfText(text, ""); ok {
			t.Errorf("DeckOfText(%q) reads a deck", text)
		}
	}
	if got := DeckText(LootDeck, DeckTextBottom); got != "loot deck bottom" {
		t.Errorf("DeckText = %q: the client reads \"loot deck bottom\"", got)
	}
}
