package b2

import "testing"

func TestPestilence(t *testing.T) {
	tb := slayTable(t, "pestilence", seat("isaac"), seat("cain"))
	killWithSixes(t, tb, me, foe)
	if a, b := tb.G.Players[0].Damage, tb.G.Players[1].Damage; a != 1 || b != 1 {
		t.Errorf("damage %d and %d, want 1 each", a, b)
	}
}
