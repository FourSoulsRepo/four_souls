package b2

import "testing"

func TestMulligan(t *testing.T) {
	tb := slayTable(t, "mulligan", seat("isaac"), seat("cain"))
	killWithSixes(t, tb)
	if n := len(tb.G.Monsters); n != 2 {
		t.Errorf("%d monster slots, want 2", n)
	}
}
