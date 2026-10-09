package b2

import "testing"

func TestHanger(t *testing.T) {
	tb := slayTable(t, "hanger", seat("isaac"), seat("cain"))
	killWithSixes(t, tb)
	if n := len(tb.G.Shop); n != 1 {
		t.Errorf("%d shop slots, want 1 (the table started with none)", n)
	}
}
