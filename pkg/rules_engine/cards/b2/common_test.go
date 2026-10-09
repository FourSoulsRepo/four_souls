package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// seat is a player with a character and items.
var seat = enginetest.Seat

// table builds a b2 table; player 0 is active.
func table(t *testing.T, seats ...engine.SituationPlayer) *enginetest.Table {
	t.Helper()
	return enginetest.New(t, Set, seats...)
}

// testExtraLoot checks a character's ↷: one more loot play, on your own
// turn and on another player's turn.
func testExtraLoot(t *testing.T, character engine.CardRef) {
	t.Helper()
	tb := table(t, seat(character), seat(character))
	tb.Start(0, character, 0)
	tb.Pass(0)
	tb.Activate(1, character, 0) // in response, on another player's turn
	if got := tb.G.Turn.LootPlays; got != 2 {
		t.Errorf("active player's loot plays = %d, want 2", got)
	}
	if tb.G.Object(tb.Find(0, character)).Charged {
		t.Error("the character is still charged after ↷")
	}
	if got := tb.G.Players[1].ExtraLootPlays; got != 1 {
		t.Errorf("other player's extra loot plays = %d, want 1", got)
	}
	if rule := tb.Refused(0, character, 0); rule != "R-ABIL-07" {
		t.Errorf("second use refused by %s, want R-ABIL-07 (deactivated)", rule)
	}
}

// newGame starts a real game with fixed characters.
func newGame(t *testing.T, characters ...engine.CardRef) *engine.Game {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 7, Players: len(characters), Sets: []engine.CardSet{Set}, Characters: characters})
	if err != nil {
		t.Fatal(err)
	}
	return g
}

// hasItem reports whether player p controls an item with this card.
func hasItem(g *engine.Game, p int, card engine.CardRef) bool {
	for _, id := range g.Players[p].InPlay {
		if o := g.Object(id); o.Card == card && o.Role == engine.RoleItem {
			return true
		}
	}
	return false
}

// monsterTable builds a b2 table with the given monsters in slots.
func monsterTable(t *testing.T, monsters []engine.CardRef, seats ...engine.SituationPlayer) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{Players: seats, Monsters: monsters}, Set)
}

// lootTable: player 0 (Isaac) holds the cards, player 1 is Cain; a fly
// and a leech are in the monster slots.
func lootTable(t *testing.T, hand ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			{Character: "isaac", Hand: hand},
			seat("cain"),
		},
		Monsters: []engine.CardRef{"fly", "leech"},
	}, Set)
}

const (
	me  = "player 1 (isaac)"
	foe = "player 2 (cain)"
)

// inDiscard reports whether the loot discard holds the card.
func inDiscard(g *engine.Game, card engine.CardRef) bool {
	for _, id := range g.Discards[engine.LootDeck] {
		if g.Object(id).Card == card {
			return true
		}
	}
	return false
}

// testGain plays a coin card and checks the cents.
func testGain(t *testing.T, card engine.CardRef, n int) {
	t.Helper()
	tb := lootTable(t, card)
	tb.Play(0, card)
	if got := tb.G.Players[0].Cents; got != n {
		t.Errorf("cents = %d, want %d", got, n)
	}
	if !inDiscard(tb.G, card) {
		t.Error("the loot card is not in the loot discard (R-CARD-07)")
	}
}

// testOneOnTop plays a "look at the top 5, put 1 on top, the rest on the
// bottom" card, picking the third card by name.
func testOneOnTop(t *testing.T, card engine.CardRef, d engine.DeckKind) {
	t.Helper()
	tb := lootTable(t, card)
	top := tb.G.DeckTop(d, 5)
	name := tb.G.Object(top[2]).Card
	kept := -1 // options are names: the first card with the name is picked
	for i, id := range top {
		if kept < 0 && tb.G.Object(id).Card == name {
			kept = i
		}
	}
	tb.Play(0, card, string(name))
	if got := tb.G.DeckTop(d, 1)[0]; got != top[kept] {
		t.Errorf("top is %s, want the chosen %s", tb.G.Object(got).Card, name)
	}
	bottom := tb.G.Decks[d][:4]
	for i, id := range top {
		if i != kept && !slices.Contains(bottom, id) {
			t.Errorf("%s is not on the bottom", tb.G.Object(id).Card)
		}
	}
}

// itemTable: player 0 (Isaac) has the items and the hand; player 1 is
// Cain; a fly and a leech are in the monster slots.
func itemTable(t *testing.T, items []engine.CardRef, hand ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			{Character: "isaac", Items: items, Hand: hand},
			seat("cain"),
		},
		Monsters: []engine.CardRef{"fly", "leech"},
	}, Set)
}

// items is a shorthand for a list of cards.
func items(cards ...engine.CardRef) []engine.CardRef { return cards }
