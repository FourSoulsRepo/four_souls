package b2

import "testing"

func TestSecretRoom(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.G.ForceRolls(4)
	tb.RevealFromDeck("secret_room")
	if c := tb.G.Players[0].Cents; c != 7 {
		t.Errorf("cents = %d, want 7", c)
	}
}
