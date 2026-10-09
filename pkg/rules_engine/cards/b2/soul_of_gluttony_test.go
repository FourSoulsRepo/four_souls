package b2

import "testing"

func TestSoulOfGluttony(t *testing.T) {
	g := bonusGame(t, "soul_of_gluttony")
	for range 10 {
		g.Players[0].Hand = append(g.Players[0].Hand, g.Players[1].Hand...)
		g.Players[1].Hand = nil
	}
	if len(g.Players[0].Hand) < 10 {
		t.Skip("not enough cards")
	}
	runSome(t, g)
	if s := g.SoulValue(0); s != 1 {
		t.Errorf("soul value = %d, want 1", s)
	}
}
