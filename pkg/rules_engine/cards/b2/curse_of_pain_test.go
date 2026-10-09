package b2

import "testing"

func TestCurseOfPain(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.RevealFromDeck("curse_of_pain", foe)
	if !hasCurse(tb.G, 1, "curse_of_pain") {
		t.Fatal("Cain did not get the curse")
	}
	tb.EndTurn()
	if d := tb.G.Players[1].Damage; d != 1 {
		t.Errorf("Cain's damage = %d, want 1", d)
	}
}
