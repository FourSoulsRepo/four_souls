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

// table builds a table; player 0 is active.
func table(t *testing.T, seats ...engine.SituationPlayer) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{Players: seats}, sets...)
}

// setup builds a table from a full setup.
func setup(t *testing.T, s engine.SituationSetup) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, s, sets...)
}

// lootTable: player 0 (Isaac) holds the cards, player 1 is Cain; a fly
// (DC 2) and a leech (DC 4) are in the monster slots.
func lootTable(t *testing.T, hand ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: hand}, seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
}

// inDiscard reports whether the deck's discard holds the card.
func inDiscard(g *engine.Game, d engine.DeckKind, card engine.CardRef) bool {
	for _, id := range g.Discards[d] {
		if g.Object(id).Card == card {
			return true
		}
	}
	return false
}

const (
	me  = "player 1 (isaac)"
	foe = "player 2 (cain)"
)

// newGame starts a real game with fixed characters.
func newGame(t *testing.T, characters ...engine.CardRef) *engine.Game {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 7, Players: len(characters), Sets: sets, Characters: characters})
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

// testCharacter checks a character's ↷ (one more loot play) and its
// starting item.
func testCharacter(t *testing.T, character, item engine.CardRef) {
	t.Helper()
	tb := table(t, seat(character), seat("cain"))
	tb.Activate(0, character, 0)
	if got := tb.G.Turn.LootPlays; got != 2 {
		t.Errorf("loot plays = %d, want 2", got)
	}
	if rule := tb.Refused(0, character, 0); rule != "R-ABIL-07" {
		t.Errorf("second use refused by %s, want R-ABIL-07 (deactivated)", rule)
	}
	g := newGame(t, character, "cain")
	p := 0
	if g.Object(g.Players[0].Character).Card != character {
		p = 1
	}
	if !hasItem(g, p, item) {
		t.Errorf("%s has no %s", character, item)
	}
}

// settleAny settles, answering every question with its first option.
func settleAny(tb *enginetest.Table) {
	tb.T.Helper()
	for range 100 {
		w := tb.G.Prompt()
		if w.Kind == engine.PromptChoose {
			tb.Choose(w.Options[0])
			continue
		}
		if len(tb.G.Stack) == 0 && w.Kind == engine.PromptPriority {
			return
		}
		tb.Pass(int(w.Player))
	}
	tb.T.Fatal("settleAny: the stack never emptied")
}

// top returns the card on top of monster slot i.
func top(tb *enginetest.Table, i int) engine.CardRef {
	id, _ := tb.G.Monsters[i].TopOf()
	return tb.G.Object(id).Card
}

// topID returns the object on top of monster slot i.
func topID(tb *enginetest.Table, i int) engine.ObjectID {
	id, _ := tb.G.Monsters[i].TopOf()
	return id
}

// attackIntent declares an attack for the active player.
func attackIntent(tb *enginetest.Table) engine.Intent {
	return engine.Intent{Player: tb.G.Turn.Active, Kind: engine.IntentAttack}
}
