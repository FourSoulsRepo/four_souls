package fsp2

import "testing"

func TestPolyphemus(t *testing.T) {
	tb := armed(t, items("polyphemus"), "fatty", "monstro")
	if atk := tb.G.PlayerATK(0); atk != 3 {
		t.Errorf("ATK = %d, want 1 + 2", atk)
	}
	if dc := tb.G.Evasion(topID(tb, 0)); dc != 3 {
		t.Errorf("fatty DC = %d on your turn, want 2 + 1", dc)
	}
}
