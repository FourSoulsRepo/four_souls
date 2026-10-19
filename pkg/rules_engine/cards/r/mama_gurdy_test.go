package r

import "testing"

func TestMamaGurdy(t *testing.T) {
	tb := fight(t, "mama_gurdy")
	hurt(tb, 0, 1)
	tb.Attack("mama_gurdy", 6, 4)
	if tb.G.SoulValue(0) != 1 || len(tb.G.Players[0].Hand) != 4 {
		t.Errorf("a soul and loot 4 from the roll; hand %d", len(tb.G.Players[0].Hand))
	}
}
