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

// hasItem reports whether player p controls an item with this card.
func hasItem(g *engine.Game, p int, card engine.CardRef) bool {
	for _, id := range g.Players[p].InPlay {
		if o := g.Object(id); o.Card == card && o.Role == engine.RoleItem {
			return true
		}
	}
	return false
}

// handHas reports whether player p holds a loot card with this card.
func handHas(g *engine.Game, p int, card engine.CardRef) bool {
	for _, id := range g.Players[p].Hand {
		if g.Object(id).Card == card {
			return true
		}
	}
	return false
}

// top returns the card on top of monster slot i.
func top(tb *enginetest.Table, i int) engine.CardRef {
	id, _ := tb.G.Monsters[i].TopOf()
	return tb.G.Object(id).Card
}

// settleAny settles, answering every question with its first option.
func settleAny(tb *enginetest.Table) {
	tb.T.Helper()
	for range 200 {
		tb.ThrowPending()
		w := tb.G.Prompt()
		switch {
		case w.Kind == engine.PromptChoose:
			tb.Choose(tb.Options()[0])
			continue
		case w.Kind == engine.PromptDiscard:
			if _, err := tb.G.Apply(engine.Intent{Player: w.Player, Kind: engine.IntentDiscard, Objects: tb.G.Players[w.Player].Hand[:w.Count]}); err != nil {
				tb.T.Fatal(err)
			}
			continue
		case len(tb.G.Stack) == 0 && w.Kind == engine.PromptPriority && !tb.G.Attack.On:
			return
		case w.Kind == engine.PromptPriority:
			tb.Pass(int(w.Player))
			continue
		}
		return
	}
	tb.T.Fatal("settleAny: the stack never emptied")
}

// hurt leaves the monster in slot i with n HP.

func hurt(tb *enginetest.Table, i, n int) {
	id := tb.G.Monsters[i].Cards[len(tb.G.Monsters[i].Cards)-1]
	tb.G.Object(id).Damage = tb.G.Def(id).HP - n
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

// room returns the room on top of room slot 0.
func room(tb *enginetest.Table) engine.ObjectID {
	id, _ := tb.G.Rooms[0].TopOf()
	return id
}

// sackFirst puts the roll result on the stack first, so the other card's
// "each time a player rolls" trigger resolves first (R-ABIL-15, 19.233).
const sackFirst = "mystery_sack: ↷: Roll- 1-2: Loot 1. 3-4: Gain 4¢."
