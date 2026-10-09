package b2

import "testing"

func TestPeep(t *testing.T) {
	tb := slayTable(t, "peep", seat("isaac"), seat("cain"))
	killWithSixes(t, tb)
	found := false
	for _, s := range tb.G.Monsters {
		if top, ok := s.TopOf(); ok && tb.G.Object(top).Card == "the_bloat" {
			found = true
		}
	}
	if !found {
		t.Error("The Bloat is not in a monster slot")
	}
}
