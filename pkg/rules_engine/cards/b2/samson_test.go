package b2

import "testing"

func TestSamson(t *testing.T) {
	testExtraLoot(t, "samson")
	g := newGame(t, "samson", "isaac")
	if !hasItem(g, 0, samson.StartingItem) {
		t.Errorf("no starting item %s", samson.StartingItem)
	}
}
