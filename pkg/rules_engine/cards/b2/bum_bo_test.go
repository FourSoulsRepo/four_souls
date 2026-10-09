package b2

import "testing"

func TestBumBo(t *testing.T) {
	tb := itemTable(t, items("bum_bo"), "a_dime")
	tb.Play(0, "a_dime")
	bumbo := tb.G.Object(tb.Find(0, "bum_bo"))
	if c, l := tb.G.Players[0].Cents, bumbo.CountersOf(""); c != 0 || l != 10 {
		t.Errorf("cents %d level %d, want 0 and 10", c, l)
	}
	if a := tb.G.PlayerATK(0); a != 2 {
		t.Errorf("ATK at level 10 = %d, want 2", a)
	}
}
