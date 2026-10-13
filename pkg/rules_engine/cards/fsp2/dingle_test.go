package fsp2

import "testing"

func TestDingle(t *testing.T) {
	tb := fight(t, "dingle", "gurdy")
	// Hits on 3 deal 1; after each, Dingle rolls: 6 means 1 more damage.
	tb.G.ForceRolls(3, 6, 3, 2) // hit + 1, hit (dead at 3)
	ev := tb.Do(attackIntent(tb), "dingle")
	if n := attackRolls(ev); n != 2 {
		t.Errorf("attack rolls = %d, want 2: Dingle (HP 3) dies to 1 + 1 + 1", n)
	}
	if top(tb, 0) == "dingle" {
		t.Error("Dingle should be dead")
	}
}
