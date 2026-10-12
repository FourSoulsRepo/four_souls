package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func voidTable(t *testing.T) *enginetest.Table {
	t.Helper()
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			{Character: "apollyon", Items: []engine.CardRef{"void"}, Hand: []engine.CardRef{"bomb", "bomb", "a_penny_6"}},
			seat("cain"),
		},
		Monsters: []engine.CardRef{"fly", "leech"},
		Shop:     []engine.CardRef{"the_d6", "breakfast"},
	})
	return tb
}

func TestVoidDiscardsAndLoots(t *testing.T) {
	tb := voidTable(t)
	before := len(tb.G.Discards[engine.LootDeck])
	tb.Activate(0, "void", 0, "Discard your hand, then loot equal to the number of cards discarded.")
	if n := len(tb.G.Players[0].Hand); n != 3 {
		t.Errorf("hand = %d cards, want 3 new ones", n)
	}
	if n := len(tb.G.Discards[engine.LootDeck]) - before; n != 3 {
		t.Errorf("%d cards discarded, want 3", n)
	}
}

func TestVoidDiscardsAMonster(t *testing.T) {
	tb := voidTable(t)
	tb.Activate(0, "void", 0, "Put a monster not being attacked into discard.", "fly")
	if top, _ := tb.G.Monsters[0].TopOf(); tb.G.Object(top).Card == "fly" {
		t.Error("the fly is still in its slot")
	}
	found := false
	for _, id := range tb.G.Discards[engine.MonsterDeck] {
		found = found || tb.G.Object(id).Card == "fly"
	}
	if !found {
		t.Error("the fly is not in the monster discard")
	}
}

func TestVoidDiscardsAShopItem(t *testing.T) {
	tb := voidTable(t)
	tb.Activate(0, "void", 0, "Put a shop item into discard.", "the_d6")
	for _, id := range tb.G.ShopItems() {
		if tb.G.Object(id).Card == "the_d6" {
			t.Error("the D6 is still in the shop")
		}
	}
	if len(tb.G.ShopItems()) != 2 {
		t.Error("the shop slot was not refilled")
	}
}
