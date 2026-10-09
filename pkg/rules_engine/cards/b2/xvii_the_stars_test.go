package b2

import "testing"

func TestXVIITheStars(t *testing.T) {
	tb := lootTable(t, "xvii_the_stars")
	tb.Play(0, "xvii_the_stars")
	if n := len(tb.G.Players[0].InPlay); n != 1 {
		t.Errorf("%d items in play, want 1", n)
	}
}
