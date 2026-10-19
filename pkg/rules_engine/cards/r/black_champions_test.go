package r

import "testing"

func TestBlackChampions(t *testing.T) {
	tb := roomTable(t, "black_champions")
	if got := tb.G.MonsterATK(tb.G.Monsters[0].Cards[0]); got != 2 {
		t.Errorf("fly ATK = %d, want 2", got)
	}
}
