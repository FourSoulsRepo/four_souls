package b2

import "testing"

func TestSoulOfGreed(t *testing.T) {
	g := bonusGame(t, "soul_of_greed")
	g.Players[1].Cents = 25
	runSome(t, g)
	if s := g.SoulValue(1); s != 1 {
		t.Errorf("soul value = %d, want 1", s)
	}
}
