package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBeggar(t *testing.T) {
	tb := roomTable(t, "beggar")
	tb.G.Players[0].Cents = 6
	tb.G.ForceRolls(6)
	useRoom(tb, 0, 0)
	if len(itemsOf(tb.G, 0)) != 1 || !inDiscard(tb.G, engine.RoomDeck, "beggar") {
		t.Error("a 6 gains a treasure and puts the beggar into discard")
	}
}
