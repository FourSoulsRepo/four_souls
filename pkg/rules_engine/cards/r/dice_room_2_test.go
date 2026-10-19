package r

import "testing"

func TestDiceRoom2(t *testing.T) {
	tb := lootRoomTable(t, "", "a_penny")
	enterRoom(tb, "dice_room_2")
	if len(tb.G.Players[0].Hand) != 3 || len(tb.G.Players[1].Hand) != 3 {
		t.Error("each player discards their hand and loots 3")
	}
}
