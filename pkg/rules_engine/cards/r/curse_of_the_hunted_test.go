package r

import "testing"

func TestCurseOfTheHunted(t *testing.T) {
	tb := fight(t, "fly")
	giveCurse(tb, 0, "curse_of_the_hunted")
	if tb.G.MonsterATK(tb.G.Monsters[0].Cards[0]) != 2 {
		t.Error("monsters have +1 ATK on his turn")
	}
}
