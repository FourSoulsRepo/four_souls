package b2

import "testing"

func TestCursedGaper(t *testing.T) {
	tb := slayTable(t, "cursed_gaper", seat("isaac", "mystery_sack"), seat("cain"))
	m, _ := tb.G.Monsters[0].TopOf()
	tb.G.ForceRolls(4)
	tb.Activate(0, "mystery_sack", 0)
	if a := tb.G.MonsterATK(m); a != 2 {
		t.Errorf("monster ATK = %d, want 2", a)
	}
}
