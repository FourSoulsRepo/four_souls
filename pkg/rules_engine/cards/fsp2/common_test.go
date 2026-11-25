package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// Four Souls+ cards are played with the Base Game.
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

// attackIntent declares an attack for the active player.
func attackIntent(tb *enginetest.Table) engine.Intent {
	return engine.Intent{Player: tb.G.Turn.Active, Kind: engine.IntentAttack}
}

// fight builds a table where Isaac (HP 2, ATK 1) faces the monsters.
func fight(t *testing.T, monsters ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Monsters: monsters,
	})
}

// armed builds a table where Isaac has the items and faces the monsters.
func armed(t *testing.T, items []engine.CardRef, monsters ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: items}, seat("cain")},
		Monsters: monsters,
	})
}

// items is a list of cards.
func items(cards ...engine.CardRef) []engine.CardRef { return cards }

// viewDice are the die numbers the spectator view shows on player p's
// card ref (15.67).
func viewDice(tb *enginetest.Table, p int, ref engine.CardRef) []int {
	tb.T.Helper()
	for _, c := range tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Players[p].InPlay {
		if c.Card == ref {
			return c.Dice
		}
	}
	tb.T.Fatalf("%s is not in play", ref)
	return nil
}
