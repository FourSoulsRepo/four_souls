package r

import "testing"

func TestEternalChest(t *testing.T) {
	tb := lootRoomTable(t, "eternal_chest", "a_penny")
	tb.G.ForceRolls(4)
	useRoom(tb, 0, 0, "a_penny")
	if got := tb.G.Players[0].Cents; got != 7 {
		t.Errorf("cents = %d, want 7", got)
	}
}
