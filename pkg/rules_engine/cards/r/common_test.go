package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// Requiem cards are played with the Base Game; its reprints need it.
// It is filled in init, after backs.go added the back faces to Set.
var sets []engine.CardSet

func init() { sets = []engine.CardSet{b2.Set, Set} }

// seat is a player with a character and items.
var seat = enginetest.Seat

// setup builds a table from a full setup. Refills and treasure gains
// draw plain cards first, so a test never depends on the deal: Clotty
// has no abilities, Book of Sin only a ↷ one.
func setup(t *testing.T, s engine.SituationSetup) *enginetest.Table {
	t.Helper()
	tb := enginetest.NewSetup(t, s, sets...)
	for range 6 {
		tb.G.PutOnDeck(engine.MonsterDeck, "clotty")
		tb.G.PutOnDeck(engine.TreasureDeck, "book_of_sin")
	}
	return tb
}

// foe is Cain's player label.
const foe = "player 2 (cain)"

// hasItem reports whether player p controls an item with this card.
func hasItem(g *engine.Game, p int, card engine.CardRef) bool {
	for _, id := range g.Players[p].InPlay {
		if o := g.Object(id); o.Card == card && o.Role == engine.RoleItem {
			return true
		}
	}
	return false
}

// testStartingItem checks that a real game gives the character's player
// its starting item (R-SETUP-07). Its ↷ is a situation file each.
func testStartingItem(t *testing.T, character, item engine.CardRef) {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 7, Players: 3, Sets: sets, Characters: []engine.CardRef{character, "cain", "isaac"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, pl := range g.Players {
		if g.Object(pl.Character).Card != character {
			continue
		}
		for _, id := range pl.InPlay {
			if g.Object(id).Card == item {
				return
			}
		}
	}
	t.Errorf("%s has no %s", character, item)
}

// roomTable: Isaac (active) and Cain, a fly and a leech, and the room.
func roomTable(t *testing.T, room engine.CardRef) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
		Rooms:    []engine.CardRef{room},
	})
}
