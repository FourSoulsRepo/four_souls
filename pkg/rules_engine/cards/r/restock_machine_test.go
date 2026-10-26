package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestRestockMachine(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
		Shop:     []engine.CardRef{"jawbone", "spoon_bender"},
		Rooms:    []engine.CardRef{"restock_machine"},
	})
	tb.G.Players[0].Cents = 2
	useRoom(tb, 0, 0, "card.restock_machine.monsters")
	if !inDiscard(tb.G, engine.MonsterDeck, "fly") || inDiscard(tb.G, engine.TreasureDeck, "jawbone") {
		t.Error("the monsters go to discard, the shop stays")
	}
}
