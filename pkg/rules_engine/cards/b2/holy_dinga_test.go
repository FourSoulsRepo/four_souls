package b2

import "testing"

func TestHolyDinga(t *testing.T) {
	tb := slayTable(t, "holy_dinga", seat("isaac", "mystery_sack"), seat("cain"))
	tb.G.Players[0].Damage = 1
	tb.G.ForceRolls(6)
	tb.Activate(0, "mystery_sack", 0)
	if d := tb.G.Players[0].Damage; d != 0 {
		t.Errorf("damage = %d, want 0", d)
	}
}
