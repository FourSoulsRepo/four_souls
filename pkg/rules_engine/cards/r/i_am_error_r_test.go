package r

import "testing"

func TestIAmErrorR(t *testing.T) {
	tb := lootRoomTable(t, "i_am_error_r", "a_penny")
	tb.G.ForceRolls(5)
	tb.EndTurn()
	if len(tb.G.Players[1].Hand) != 4 {
		t.Errorf("a 5: Cain discards his hand and loots 3, then loots in his loot step: hand %d", len(tb.G.Players[1].Hand))
	}
}
