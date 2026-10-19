package r

import "testing"

func TestKeyMaster(t *testing.T) {
	tb := lootRoomTable(t, "key_master", "a_penny")
	tb.G.ForceRolls(6)
	useRoom(tb, 0, 0, "a_penny")
	if len(itemsOf(tb.G, 0)) != 1 {
		t.Error("a 6 gains a treasure")
	}
}
