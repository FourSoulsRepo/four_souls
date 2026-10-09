package b2

import "testing"

func TestBoomFly(t *testing.T) {
	tb := slayTable(t, "boom_fly", seat("isaac"), seat("cain"))
	killWithSixes(t, tb)
	if a, b := tb.G.Players[0].Damage, tb.G.Players[1].Damage; a != 1 || b != 1 {
		t.Errorf("damage %d and %d, want 1 each", a, b)
	}
}
