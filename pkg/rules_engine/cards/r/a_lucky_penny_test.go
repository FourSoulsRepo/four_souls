package r

import "testing"

func TestALuckyPenny(t *testing.T) {
	tb := lootTable(t, "a_lucky_penny", "a_lucky_penny")
	tb.Play(0, "a_lucky_penny")
	if got := tb.G.Players[0].Cents; got != 1 {
		t.Fatalf("cents = %d, want 1", got)
	}
	tb.Attack("fly", 6) // the fly gives 1¢, doubled
	if got := tb.G.Players[0].Cents; got != 3 {
		t.Errorf("cents after the fly = %d, want 1 + 2", got)
	}
}

func TestALuckyPennyDoublesOnce(t *testing.T) {
	tb := lootTable(t, "a_lucky_penny", "a_lucky_penny")
	tb.G.Turn.LootPlays = 2
	tb.Play(0, "a_lucky_penny")
	tb.Play(0, "a_lucky_penny")
	tb.Attack("fly", 6)
	if got := tb.G.Players[0].Cents; got != 4 {
		t.Errorf("cents = %d, want 2 + 2: rewards double once (R-ABIL-36)", got)
	}
}
