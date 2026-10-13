package g2

import "testing"

func TestRoundWorm(t *testing.T) {
	tb := fight(t, "round_worm", "fly")
	tb.G.ForceRolls(2, 4) // miss: the worm hits; the next roll gets +1: 5 hits
	ev := tb.Do(attackIntent(tb), "round_worm")
	if n := attackRolls(ev); n != 2 {
		t.Errorf("attack rolls = %d, want 2", n)
	}
}
