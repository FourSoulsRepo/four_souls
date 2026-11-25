package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// Gold Box cards are played with the Base Game.
var sets = []engine.CardSet{b2.Set, Set}

// seat is a player with a character and items.
var seat = enginetest.Seat

// setup builds a table from a full setup.
func setup(t *testing.T, s engine.SituationSetup) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, s, sets...)
}

// me is Isaac's label in a two-player table.
const me = "player 1 (isaac)"

// newGame starts a real game with fixed characters.
func newGame(t *testing.T, characters ...engine.CardRef) *engine.Game {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 7, Players: len(characters), Sets: sets, Characters: characters})
	if err != nil {
		t.Fatal(err)
	}
	enginetest.ThrowFirstRoll(t, g)
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

// testStartingItem checks a character's starting item in a real game
// start. Its ↷ is in testdata/situations/<character>/.
func testStartingItem(t *testing.T, character, item engine.CardRef) {
	t.Helper()
	g := newGame(t, character, "cain")
	p := 0
	if g.Object(g.Players[0].Character).Card != character {
		p = 1
	}
	if !hasItem(g, p, item) {
		t.Errorf("%s has no %s", character, item)
	}
}

// top returns the card on top of monster slot i.
func top(tb *enginetest.Table, i int) engine.CardRef {
	id, _ := tb.G.Monsters[i].TopOf()
	return tb.G.Object(id).Card
}
