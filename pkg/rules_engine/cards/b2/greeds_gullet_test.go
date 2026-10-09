package b2

import "testing"

func TestGreedsGullet(t *testing.T) {
	tb := itemTable(t, items("greeds_gullet"), "xiii_death")
	tb.Play(0, "xiii_death", me, "greeds_gullet")
	// 8¢ before the penalty; then the penalty takes 1¢ and the gullet.
	if c := tb.G.Players[0].Cents; c != 7 {
		t.Errorf("cents = %d, want 7", c)
	}
}
