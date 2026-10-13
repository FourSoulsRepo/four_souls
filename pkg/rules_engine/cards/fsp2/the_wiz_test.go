package fsp2

import "testing"

func TestTheWiz(t *testing.T) {
	tb := armed(t, items("the_wiz"), "fly", "gurdy")
	if dc := tb.G.Evasion(topID(tb, 0)); dc != 3 {
		t.Errorf("fly DC = %d on your turn, want 3", dc)
	}
	tb.G.ForceRolls(3)
	tb.Do(attackIntent(tb), "fly", "gurdy")
	if d := tb.G.Object(topID(tb, 1)).Damage; d != 1 {
		t.Errorf("Gurdy damage = %d, want 1", d)
	}
}
