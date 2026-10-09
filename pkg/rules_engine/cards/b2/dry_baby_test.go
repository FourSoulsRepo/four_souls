package b2

import "testing"

func TestDryBaby(t *testing.T) {
	tb := monsterTable(t, items("leech"), seat("isaac", "dry_baby"), seat("cain"))
	tb.Attack("leech", 1, 6) // 2 damage reduced to 1
	if d := tb.G.Players[0].Damage; d != 1 {
		t.Errorf("damage = %d, want 1", d)
	}
}
