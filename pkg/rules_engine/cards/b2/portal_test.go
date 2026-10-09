package b2

import "testing"

func TestPortal(t *testing.T) {
	tb := slayTable(t, "portal", seat("isaac"), seat("cain"))
	killWithSixes(t, tb)
	if n := tb.G.Turn.MustAttacks; n != 1 {
		t.Errorf("must attacks = %d, want 1", n)
	}
}
