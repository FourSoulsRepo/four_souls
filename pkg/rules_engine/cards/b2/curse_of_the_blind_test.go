package b2

import "testing"

func TestCurseOfTheBlind(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.RevealFromDeck("curse_of_the_blind", foe)
	if !hasCurse(tb.G, 1, "curse_of_the_blind") {
		t.Fatal("Cain did not get the curse")
	}
	tb.EndTurn()
	fly, _ := tb.G.Monsters[0].TopOf()
	if dc := tb.G.Evasion(fly); dc != 3 {
		t.Errorf("fly DC on Cain's turn = %d, want 3", dc)
	}
}
