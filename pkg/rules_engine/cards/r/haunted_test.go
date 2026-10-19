package r

import "testing"

func TestHaunted(t *testing.T) {
	tb := fightRoom(t, "haunted", "chub")
	hurt(tb, 0, 1)
	tb.Attack("chub", 6)
	if tb.G.Players[0].SkipTurns != 1 {
		t.Error("Isaac gained a soul: he skips his next turn")
	}
}
