package b2

import "testing"

func TestBlueBaby(t *testing.T) {
	testExtraLoot(t, "blue_baby")
	g := newGame(t, "blue_baby", "isaac")
	if !hasItem(g, 0, blueBaby.StartingItem) {
		t.Errorf("no starting item %s", blueBaby.StartingItem)
	}
}
