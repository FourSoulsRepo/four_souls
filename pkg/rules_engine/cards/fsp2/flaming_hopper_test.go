package fsp2

import "testing"

func TestFlamingHopper(t *testing.T) {
	tb := armed(t, items("lard"), "flaming_hopper", "fly")
	tb.G.ForceRolls(6, 4) // a 6 deals nothing; a 4 kills it (HP 1)
	ev := tb.Do(attackIntent(tb), "flaming_hopper")
	if n := attackRolls(ev); n != 2 {
		t.Errorf("attack rolls = %d, want 2", n)
	}
}
