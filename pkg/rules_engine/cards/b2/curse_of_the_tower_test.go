package b2

import "testing"

func TestCurseOfTheTower(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac", "curse_of_the_tower"), seat("cain"))
	tb.Attack("fly", 1, 2, 6) // miss: the curse rolls 2, Cain takes 1; then kill
	if d := tb.G.Players[1].Damage; d != 1 {
		t.Errorf("Cain's damage = %d, want 1", d)
	}
}
