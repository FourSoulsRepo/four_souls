package b2

import "testing"

func TestOTheFool(t *testing.T) {
	tb := lootTable(t, "o_the_fool")
	tb.Play(0, "o_the_fool")
	if tb.G.Turn.Active != 1 && !tb.G.Turn.EndDeclared {
		t.Error("the turn did not end")
	}
}
