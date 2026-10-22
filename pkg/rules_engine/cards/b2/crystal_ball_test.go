package b2

import "testing"

func TestCrystalBall(t *testing.T) {
	tb := itemTable(t, items("crystal_ball", "mystery_sack"))
	tb.Activate(0, "crystal_ball", 0, "4")
	if d := viewDice(tb, 0, "crystal_ball"); len(d) != 1 || d[0] != 4 {
		t.Errorf("view dice = %v, want [4]", d)
	}
	tb.G.ForceRolls(4)
	tb.Activate(0, "mystery_sack", 0)
	if h := len(tb.G.Players[0].Hand); h != 3 {
		t.Errorf("hand = %d, want 3", h)
	}
	if tb.G.Object(tb.Find(0, "crystal_ball")).Chosen != 0 {
		t.Error("the guess was not used up")
	}
	if d := viewDice(tb, 0, "crystal_ball"); d != nil {
		t.Errorf("view dice = %v after the roll, want none", d)
	}
}
