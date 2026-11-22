package b2

import "testing"

func TestEve(t *testing.T) {
	testExtraLootOnAnotherTurn(t, "eve")
	g := newGame(t, "eve", "isaac")
	if !hasItem(g, 0, eve.StartingItem) {
		t.Errorf("no starting item %s", eve.StartingItem)
	}
}
