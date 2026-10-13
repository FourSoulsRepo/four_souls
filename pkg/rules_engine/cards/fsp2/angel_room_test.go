package fsp2

import "testing"

func TestAngelRoom(t *testing.T) {
	for _, c := range []struct{ roll, items, hand int }{{1, 2, 0}, {3, 1, 0}, {5, 0, 2}} {
		tb := fight(t, "fatty", "monstro")
		tb.G.ForceRolls(c.roll)
		tb.RevealFromDeck("angel_room")
		if n := len(tb.G.Players[0].InPlay); n != c.items {
			t.Errorf("roll %d: items = %d, want %d", c.roll, n, c.items)
		}
		if n := len(tb.G.Players[0].Hand); n != c.hand {
			t.Errorf("roll %d: hand = %d, want %d", c.roll, n, c.hand)
		}
	}
}
