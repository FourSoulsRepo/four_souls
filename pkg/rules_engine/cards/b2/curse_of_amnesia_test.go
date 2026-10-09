package b2

import "testing"

func TestCurseOfAmnesia(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.RevealFromDeck("curse_of_amnesia", foe)
	if !hasCurse(tb.G, 1, "curse_of_amnesia") {
		t.Fatal("Cain did not get the curse")
	}
	tb.G.Players[1].Hand = nil
	tb.EndTurn() // Cain loots 1 in the loot step
	first := string(tb.G.Object(tb.G.Players[1].Hand[0]).Card)
	tb.EndTurn(first) // then discards it; the second discard finds nothing
	if h := len(tb.G.Players[1].Hand); h != 0 {
		t.Errorf("Cain's hand %d, want 0", h)
	}
}
