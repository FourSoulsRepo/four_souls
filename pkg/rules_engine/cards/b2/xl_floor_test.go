package b2

import "testing"

func TestXLFloor(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.RevealFromDeck("xl_floor")
	if n := len(tb.G.Monsters); n != 2 {
		t.Errorf("%d monster slots, want 2", n)
	}
	if _, ok := tb.G.Monsters[1].TopOf(); !ok {
		t.Error("the new slot is empty")
	}
}
