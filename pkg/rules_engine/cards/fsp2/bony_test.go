package fsp2

import "testing"

// The +2 ATK comes after the 1's own miss is on the stack, so it counts
// from the next miss (R-ATK-13: the miss's damage is set as it goes on the
// stack).
func TestBony(t *testing.T) {
	tb := armed(t, items("lard", "lard"), "bony", "fly") // Isaac HP 6
	tb.G.ForceRolls(1, 1, 3, 3)                          // misses for 1, then 3; two hits kill Bony
	lard := "lard: Each time you take damage, discard a loot card."
	tb.Do(attackIntent(tb), "bony", lard, lard) // the two Lards' order, once per miss
	if d := tb.G.Players[0].Damage; d != 4 {
		t.Errorf("damage = %d, want 1 + 3", d)
	}
}
