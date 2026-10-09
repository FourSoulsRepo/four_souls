package b2

import "testing"

func TestChampionBelt(t *testing.T) {
	tb := itemTable(t, items("champion_belt"))
	if n := len(tb.G.Allowed(0)); n == 0 {
		t.Fatal("nothing allowed")
	}
	tb.Attack("fly", 6)
	tb.Attack("leech", 6) // a second attack
	if tb.G.Turn.BonusAttacksUsed != 1 {
		t.Errorf("bonus attacks used = %d, want 1", tb.G.Turn.BonusAttacksUsed)
	}
}
