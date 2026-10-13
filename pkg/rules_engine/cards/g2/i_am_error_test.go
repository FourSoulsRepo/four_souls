package g2

import "testing"

func TestIAmError(t *testing.T) {
	for _, c := range []struct{ roll, hand, must int }{{1, 0, 1}, {3, 3, 0}} {
		tb := fight(t, "fatty", "monstro")
		tb.G.ForceRolls(c.roll)
		tb.RevealFromDeck("i_am_error")
		if n := len(tb.G.Players[0].Hand); n != c.hand {
			t.Errorf("roll %d: hand = %d, want %d", c.roll, n, c.hand)
		}
		if tb.G.Turn.MustAttackDeck != c.must {
			t.Errorf("roll %d: must attack the deck %d times, want %d", c.roll, tb.G.Turn.MustAttackDeck, c.must)
		}
	}
	tb := fight(t, "fatty", "monstro")
	tb.G.ForceRolls(6)
	tb.RevealFromDeck("i_am_error")
	if len(tb.G.Players[0].InPlay) != 1 {
		t.Error("5-6: gain a treasure")
	}
}
