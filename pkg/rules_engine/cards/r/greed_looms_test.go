package r

import "testing"

func TestGreedLooms(t *testing.T) {
	tb := lootRoomTable(t, "greed_looms", "a_penny", "a_dime")
	tb.EndTurn()
	if len(tb.G.Players[0].Hand) != 0 {
		t.Error("Isaac bought nothing: he discards his hand")
	}
}
