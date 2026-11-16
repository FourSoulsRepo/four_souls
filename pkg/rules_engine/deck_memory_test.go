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
