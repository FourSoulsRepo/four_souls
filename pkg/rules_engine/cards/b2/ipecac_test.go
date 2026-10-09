package b2

import "testing"

func TestIpecac(t *testing.T) {
	tb := monsterTable(t, items("conjoined_fatty"), seat("isaac", "ipecac"), seat("cain"))
	tb.Attack("conjoined_fatty", 6, 6)
	if d := tb.G.Players[1].Damage; d != 2 {
		t.Errorf("Cain's damage = %d, want 2 (two rolls of 6)", d)
	}
}
