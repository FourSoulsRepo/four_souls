package r

import "testing"

func TestNull(t *testing.T) {
	tb := fight(t, "null")
	tb.Attack("null", 6, 6, 6, 6, 6) // the 2nd and 4th hits deal nothing
	if top(tb, 0) == "null" {
		t.Fatal("the 1st, 3rd and 5th hits kill it (3 HP)")
	}
	if tb.G.Turn.AttackRolls != 5 {
		t.Errorf("attack rolls = %d, want 5", tb.G.Turn.AttackRolls)
	}
}
