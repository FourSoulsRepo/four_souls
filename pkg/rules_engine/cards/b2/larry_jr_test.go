package b2

import "testing"

func TestLarryJr(t *testing.T) {
	tb := slayTable(t, "larry_jr", seat("isaac"), seat("cain"))
	l, _ := tb.G.Monsters[0].TopOf()
	tb.G.Object(l).Damage = 2
	if dc := tb.G.Evasion(l); dc != 4 {
		t.Errorf("DC at 2 HP = %d, want 4", dc)
	}
}
