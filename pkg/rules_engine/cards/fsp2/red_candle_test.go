package fsp2

import "testing"

func TestRedCandle(t *testing.T) {
	tb := armed(t, items("red_candle", "mystery_sack"), "fly", "fatty")
	tb.Activate(0, "red_candle", 0, "3")
	if d := viewDice(tb, 0, "red_candle"); len(d) != 1 || d[0] != 3 {
		t.Errorf("view dice = %v, want [3]", d)
	}
	tb.G.ForceRolls(3)
	tb.Activate(0, "mystery_sack", 0, "fly")
	if top(tb, 0) == "fly" {
		t.Error("a 3 was rolled: 1 damage kills the fly")
	}
	tb.EndTurn()
	if len(tb.G.Lingering) != 0 {
		t.Error("it lasts till end of turn")
	}
	if d := viewDice(tb, 0, "red_candle"); d != nil {
		t.Errorf("view dice = %v after the turn, want none", d)
	}
}
