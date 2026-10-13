package g2

import "testing"

func TestPolycephalus(t *testing.T) {
	tb := fight(t, "polycephalus", "fly")
	tb.G.ForceRolls(4, 4, 3, 3) // hit; next has -2: miss; hit; hit (dies at 3)
	ev := tb.Do(attackIntent(tb), "polycephalus")
	if n := attackRolls(ev); n != 4 {
		t.Errorf("attack rolls = %d, want 4", n)
	}
}
