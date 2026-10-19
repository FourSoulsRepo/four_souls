package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestFortuneTeller(t *testing.T) {
	tb := roomTable(t, "fortune_teller")
	n := len(tb.G.Players[0].Hand)
	tb.G.Players[0].Cents = 3
	tb.G.ForceRolls(6)
	useRoom(tb, 0, 0)
	if got := len(tb.G.Players[0].Hand) - n; got != 6 || !inDiscard(tb.G, engine.RoomDeck, "fortune_teller") {
		t.Errorf("looted %d; a 6 loots 6 and puts this into discard", got)
	}
}
