package g2

import "testing"

func TestLibraryCard(t *testing.T) {
	tb := table(t, seat("isaac", "library_card"), seat("cain"))
	tb.Activate(0, "library_card", 0)
	if got := tb.G.Turn.LootPlays; got != 2 {
		t.Errorf("loot plays = %d, want 2", got)
	}
}
