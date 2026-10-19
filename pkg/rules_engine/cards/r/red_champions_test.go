package r

import "testing"

func TestRedChampions(t *testing.T) {
	tb := roomTable(t, "red_champions")
	if got := tb.G.HP(tb.G.Monsters[0].Cards[0]); got != 2 {
		t.Errorf("fly HP = %d, want 2", got)
	}
}
