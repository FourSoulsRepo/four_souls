package b2

import "testing"

func TestFamine(t *testing.T) {
	tb := slayTable(t, "famine", seat("isaac"), seat("cain"))
	killWithSixes(t, tb)
	tb.EndTurn()
	if tb.G.Turn.Active != 1 {
		t.Fatal("Cain's turn did not come")
	}
	tb.EndTurn()
	if tb.G.Turn.Active != 1 {
		t.Errorf("Isaac's turn was not skipped: player %d is active", tb.G.Turn.Active)
	}
}
