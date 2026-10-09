package b2

import "testing"

func TestGreedEvent(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.G.Players[1].Cents = 9
	tb.RevealFromDeck("greed_event", foe)
	if c := tb.G.Players[1].Cents; c != 0 {
		t.Errorf("Cain's cents = %d, want 0", c)
	}
}
