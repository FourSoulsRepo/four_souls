package r

import "testing"

func TestStickyNickelGain(t *testing.T) {
	tb := lootTable(t, "sticky_nickel")
	tb.G.ForceRolls(4)
	tb.Play(0, "sticky_nickel")
	if got := tb.G.Players[0].Cents; got != 5 {
		t.Errorf("cents = %d, want 5", got)
	}
}

func TestStickyNickelSticks(t *testing.T) {
	tb := lootTable(t, "sticky_nickel")
	tb.G.ForceRolls(2)
	tb.Play(0, "sticky_nickel")
	if !handHas(tb.G, 1, "sticky_nickel") {
		t.Error("1-3: the player to the left should hold the nickel")
	}
	if inDiscard(tb.G, 1, "sticky_nickel") {
		t.Error("the nickel is still in the loot discard")
	}
}
