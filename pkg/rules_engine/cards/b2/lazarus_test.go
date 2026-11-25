package b2

import "testing"

func TestLazarus(t *testing.T) {
	testExtraLootOnAnotherTurn(t, "lazarus")
	g := newGame(t, "lazarus", "isaac")
	if !hasItem(g, 0, lazarus.StartingItem) {
		t.Errorf("no starting item %s", lazarus.StartingItem)
	}
}
