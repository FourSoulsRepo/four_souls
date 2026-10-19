package r

import "testing"

func TestSpiderWebs(t *testing.T) {
	tb := roomTable(t, "spider_webs")
	tb.EndTurn()
	if tb.G.Object(tb.G.Players[0].Character).Charged {
		t.Error("Isaac's character is deactivated at the end of his turn")
	}
}
