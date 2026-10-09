package b2

import "testing"

func TestCurseOfGreed(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.RevealFromDeck("curse_of_greed", foe)
	if !hasCurse(tb.G, 1, "curse_of_greed") {
		t.Fatal("Cain did not get the curse")
	}
	tb.G.Players[1].Cents = 5
	tb.EndTurn()
	tb.EndTurn()
	if c := tb.G.Players[1].Cents; c != 1 {
		t.Errorf("Cain's cents = %d, want 1", c)
	}
}
