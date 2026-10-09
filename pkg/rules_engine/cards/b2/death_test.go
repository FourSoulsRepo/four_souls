package b2

import "testing"

func TestDeath(t *testing.T) {
	tb := slayTable(t, "death", seat("isaac"), seat("cain"))
	killWithSixes(t, tb, foe)
	if !tb.G.Players[1].Dead {
		t.Error("Cain is not dead")
	}
}
