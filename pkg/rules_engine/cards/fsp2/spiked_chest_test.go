package fsp2

import "testing"

func TestSpikedChest(t *testing.T) {
	tb := fight(t, "fatty", "monstro")
	tb.G.ForceRolls(3)
	tb.RevealFromDeck("spiked_chest")
	if tb.G.Players[0].Damage != 1 || len(tb.G.Players[0].Hand) != 1 {
		t.Error("3-4: take 1 damage, loot 1")
	}
}
