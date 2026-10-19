package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSlotMachine(t *testing.T) {
	tb := roomTable(t, "slot_machine")
	tb.G.Players[0].Cents = 3
	tb.G.ForceRolls(6)
	useRoom(tb, 0, 0)
	if tb.G.Players[0].Cents != 15 || !inDiscard(tb.G, engine.RoomDeck, "slot_machine") {
		t.Errorf("cents = %d, want 15; a 6 puts this into discard", tb.G.Players[0].Cents)
	}
}
