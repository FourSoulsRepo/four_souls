package b2

import "testing"

func TestIsaac(t *testing.T) {
	testExtraLoot(t, "isaac")
	g := newGame(t, "isaac", "isaac")
	if !hasItem(g, 0, isaac.StartingItem) {
		t.Errorf("no starting item %s", isaac.StartingItem)
	}
}
