package fsp2

import "testing"

func TestRoundy(t *testing.T) {
	tb := armed(t, items("lard"), "roundy", "gurdy")
	tb.G.ForceRolls(1, 6, 6, 6) // a miss: Roundy deals combat damage, then 1 to Gurdy
	tb.Do(attackIntent(tb), "roundy", "roundy: Each time this deals combat damage, it deals 1 damage to each other monster.")
	if d := tb.G.Object(topID(tb, 1)).Damage; d != 1 {
		t.Errorf("Gurdy damage = %d, want 1", d)
	}
}
