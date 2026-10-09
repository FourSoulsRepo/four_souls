package b2

import "testing"

func TestMaggy(t *testing.T) {
	testExtraLoot(t, "maggy")
	g := newGame(t, "maggy", "isaac")
	if !hasItem(g, 0, maggy.StartingItem) {
		t.Errorf("no starting item %s", maggy.StartingItem)
	}
}
