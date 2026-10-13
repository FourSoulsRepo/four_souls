package fsp2

import "testing"

func TestGoldKey(t *testing.T) {
	tb := lootTable(t, "gold_key")
	tb.Play(0, "gold_key")
	tb.Attack("fly", 6)
	if tb.G.Turn.DeckAttacks < 50 {
		t.Errorf("deck attacks left = %d, want many", tb.G.Turn.DeckAttacks)
	}
}
