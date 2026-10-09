package b2

import "testing"

func TestTheForgotten(t *testing.T) {
	testExtraLoot(t, "the_forgotten")
	g := newGame(t, "the_forgotten", "isaac")
	if !hasItem(g, 0, theForgotten.StartingItem) {
		t.Errorf("no starting item %s", theForgotten.StartingItem)
	}
}
