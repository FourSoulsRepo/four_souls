package g2

import "testing"

func TestParabite(t *testing.T) {
	tb := fight(t, "parabite", "fly")
	tb.G.ForceRolls(3, 3, 3) // hit; the next roll gets -1: miss; then hit
	ev := tb.Do(attackIntent(tb), "parabite")
	if n := attackRolls(ev); n != 3 {
		t.Errorf("attack rolls = %d, want 3", n)
	}
	if d := tb.G.Players[0].Damage; d != 1 {
		t.Errorf("damage = %d, want 1 from the missed second roll", d)
	}
}
