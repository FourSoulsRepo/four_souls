package fsp2

import "testing"

func TestHolyChest(t *testing.T) {
	tb := fight(t, "fatty", "monstro")
	tb.G.ForceRolls(6)
	tb.RevealFromDeck("holy_chest")
	if tb.G.SoulValue(0) != 1 {
		t.Error("6: the chest becomes a soul")
	}
	tb2 := fight(t, "fatty", "monstro")
	tb2.G.ForceRolls(4)
	tb2.RevealFromDeck("holy_chest")
	if c := tb2.G.Players[0].Cents; c != 7 {
		t.Errorf("3-5: cents = %d, want 7", c)
	}
}
