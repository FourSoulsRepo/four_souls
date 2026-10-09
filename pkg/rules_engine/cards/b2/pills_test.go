package b2

import "testing"

func TestPills(t *testing.T) {
	// The hand keeps A Dime!!: loot 1, loot 3, or discard it.
	for roll, hand := range map[int]int{1: 2, 3: 4, 5: 0} {
		tb := lootTable(t, "pills", "a_dime")
		tb.G.ForceRolls(roll)
		tb.Play(0, "pills", "a_dime")
		if got := len(tb.G.Players[0].Hand); got != hand {
			t.Errorf("roll %d: hand %d, want %d", roll, got, hand)
		}
	}
}
