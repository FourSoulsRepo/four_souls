package r

import "testing"

func TestButterFingers(t *testing.T) {
	tb := lootRoomTable(t, "butter_fingers", "a_penny", "a_dime")
	tb.EndTurn("a_dime")
	if len(tb.G.Players[0].Hand) != 1 || !handHas(tb.G, 0, "a_penny") {
		t.Error("Isaac discards a loot card at the end of his turn")
	}
}
