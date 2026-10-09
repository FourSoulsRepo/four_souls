package b2

import "testing"

func TestBlackBony(t *testing.T) {
	tb := slayTable(t, "black_bony", seat("isaac"), seat("cain"))
	killWithSixes(t, tb)
	if d := tb.G.Players[0].Damage; d != 1 {
		t.Errorf("the killer's damage = %d, want 1", d)
	}
}
