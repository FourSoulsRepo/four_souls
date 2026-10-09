package b2

import "testing"

func TestMulliboom(t *testing.T) {
	tb := slayTable(t, "mulliboom", seat("isaac"), seat("cain"))
	killWithSixes(t, tb, foe)
	if !tb.G.Players[1].Dead {
		t.Error("3 damage did not kill Cain")
	}
}
