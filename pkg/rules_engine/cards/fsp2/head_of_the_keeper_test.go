package fsp2

import "testing"

func TestHeadOfTheKeeper(t *testing.T) {
	tb := armed(t, items("head_of_the_keeper"), "fatty", "monstro")
	tb.G.ForceRolls(2, 2, 2, 2) // four hits kill the fatty (4 HP)
	tb.Do(attackIntent(tb), "fatty")
	// 4 hits × 1¢ and the kill 2¢; the fatty's reward is a loot card.
	if c := tb.G.Players[0].Cents; c != 4+2 {
		t.Errorf("cents = %d, want 6", c)
	}
}
