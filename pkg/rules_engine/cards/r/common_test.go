package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// Requiem cards are played with the Base Game; its reprints need it.
var sets = []engine.CardSet{b2.Set, Set}

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
			tb.Choose(w.Options[0])
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

// Tests of the simple passive treasures; the others have their own files.

// table2 seats Isaac with the items and Cain.

func table2(t *testing.T, items ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", items...), seat("cain")}})
}

// Attacks on players, attackable items and helpers (R-ATK-20, R-ATK-21,
// R-MECH-04), through the cards that allow them.

func sirenTable(t *testing.T, isaacItems, cainItems []engine.CardRef) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", isaacItems...), seat("cain", cainItems...)},
		Monsters: []engine.CardRef{"the_siren", "fly"},
	})
}

// fight: Isaac (active), Cain and Samson; the monster in slot 1, a fly

// in slot 2.

func fight(t *testing.T, monster engine.CardRef, isaacItems ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", isaacItems...), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{monster, "fly"},
	})
}

// hurt leaves the monster in slot i with n HP.

func hurt(tb *enginetest.Table, i, n int) {
	id := tb.G.Monsters[i].Cards[len(tb.G.Monsters[i].Cards)-1]
	tb.G.Object(id).Damage = tb.G.Def(id).HP - n
}

// charmedGift kills the charmed monster (1 HP left) and gives its gift

// to Cain, with extra answers before that.

func charmedGift(t *testing.T, card engine.CardRef, extra ...string) *enginetest.Table {
	t.Helper()
	tb := fight(t, card)
	hurt(tb, 0, 1)
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, append(append([]string{string(card)}, extra...), foe)...)
	return tb
}

// giveCurse puts a curse card in play under player p.
func giveCurse(tb *enginetest.Table, p int, curse engine.CardRef) {
	tb.T.Helper()
	tb.G.PutOnDeck(engine.MonsterDeck, curse)
	top := tb.G.DeckTop(engine.MonsterDeck, 1)[0]
	tb.G.GiveCurseFromDeck(top, engine.PlayerID(p))
}
