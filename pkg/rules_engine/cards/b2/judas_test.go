package b2

import "testing"

func TestJudas(t *testing.T) {
	testExtraLoot(t, "judas")
	g := newGame(t, "judas", "isaac")
	if !hasItem(g, 0, judas.StartingItem) {
		t.Errorf("no starting item %s", judas.StartingItem)
	}
}
