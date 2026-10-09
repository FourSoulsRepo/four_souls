package b2

import "testing"

func TestPills2(t *testing.T) {
	for roll, cents := range map[int]int{2: 9, 4: 12, 6: 1} {
		tb := lootTable(t, "pills_2")
		tb.G.Players[0].Cents = 5
		tb.G.ForceRolls(roll)
		tb.Play(0, "pills_2")
		if got := tb.G.Players[0].Cents; got != cents {
			t.Errorf("roll %d: cents %d, want %d", roll, got, cents)
		}
	}
}
