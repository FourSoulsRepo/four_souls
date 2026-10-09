package b2

import "testing"

func TestLilith(t *testing.T) {
	testExtraLoot(t, "lilith")
	g := newGame(t, "lilith", "isaac")
	if !hasItem(g, 0, lilith.StartingItem) {
		t.Errorf("no starting item %s", lilith.StartingItem)
	}
}
