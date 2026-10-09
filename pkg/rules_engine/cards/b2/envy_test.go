package b2

import "testing"

func TestEnvy(t *testing.T) {
	tb := slayTable(t, "envy", seat("isaac"), seat("cain"))
	killWithSixes(t, tb)
	if n := tb.G.Turn.MustAttacks; n != 1 {
		t.Errorf("must attacks = %d, want 1", n)
	}
}
