package r

import "testing"

func TestEquality(t *testing.T) {
	tb := roomTable(t, "equality")
	tb.Attack("fly", 6) // the fly gives 1¢
	if tb.G.Players[0].Cents != 1 || tb.G.Players[1].Cents != 1 {
		t.Errorf("each player gains the 1¢: %d and %d", tb.G.Players[0].Cents, tb.G.Players[1].Cents)
	}
}
