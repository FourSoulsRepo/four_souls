package b2

import "testing"

func TestMomsBra(t *testing.T) {
	tb := monsterTable(t, items("leech"), seat("isaac", "moms_bra"), seat("cain"))
	tb.Activate(0, "moms_bra", 0, me)
	tb.Attack("leech", 1, 6)
	if d := tb.G.Players[0].Damage; d != 1 {
		t.Errorf("damage = %d, want 1", d)
	}
}
