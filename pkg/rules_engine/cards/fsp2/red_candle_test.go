package fsp2

import "testing"

// TestRedCandleShowsTheNumber: the view shows the chosen number till the
// end of the turn (15.67). The damage is in testdata/situations/red_candle/.
func TestRedCandleShowsTheNumber(t *testing.T) {
	tb := armed(t, items("red_candle"), "fly", "fatty")
	tb.Activate(0, "red_candle", 0, "3")
	if d := viewDice(tb, 0, "red_candle"); len(d) != 1 || d[0] != 3 {
		t.Errorf("view dice = %v, want [3]", d)
	}
	tb.EndTurn()
	if d := viewDice(tb, 0, "red_candle"); d != nil {
		t.Errorf("view dice = %v after the turn, want none", d)
	}
}
