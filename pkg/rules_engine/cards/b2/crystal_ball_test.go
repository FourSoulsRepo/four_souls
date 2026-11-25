package b2

import "testing"

// TestCrystalBallShowsTheGuess: the guess shows on the card until the
// next roll uses it up (a view). The loot 3 is a situation file:
// crystal_ball/right-guess-loots-3.json.
func TestCrystalBallShowsTheGuess(t *testing.T) {
	tb := itemTable(t, items("crystal_ball", "mystery_sack"))
	tb.Activate(0, "crystal_ball", 0, "4")
	if d := viewDice(tb, 0, "crystal_ball"); len(d) != 1 || d[0] != 4 {
		t.Errorf("view dice = %v, want [4]", d)
	}
	tb.G.ForceRolls(4)
	tb.Activate(0, "mystery_sack", 0, sackFirst)
	if tb.G.Object(tb.Find(0, "crystal_ball")).Chosen != 0 {
		t.Error("the guess was not used up")
	}
	if d := viewDice(tb, 0, "crystal_ball"); d != nil {
		t.Errorf("view dice = %v after the roll, want none", d)
	}
}
