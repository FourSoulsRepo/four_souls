package r

import "testing"

func TestWhiteChampions(t *testing.T) {
	tb := roomTable(t, "white_champions")
	if got := tb.G.Evasion(tb.G.Monsters[0].Cards[0]); got != 3 {
		t.Errorf("fly DC = %d, want 3", got)
	}
}
