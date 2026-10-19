package r

import "testing"

func TestBabyPlum(t *testing.T) {
	tb := fight(t, "baby_plum")
	tb.Attack("baby_plum", 6, 3, 6, 5, 6, 1, 6, 1)
	if tb.G.Players[1].Damage != 1 || tb.G.Players[2].Damage != 1 || tb.G.Players[0].Damage != 2 {
		t.Errorf("3-4: left (Cain), 5-6: right (Samson), 1-2: Isaac: %d %d %d", tb.G.Players[0].Damage, tb.G.Players[1].Damage, tb.G.Players[2].Damage)
	}
}
