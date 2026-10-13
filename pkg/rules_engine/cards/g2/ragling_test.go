package g2

import "testing"

func TestRagling(t *testing.T) {
	tb := fight(t, "ragling", "fly")
	tb.G.ForceRolls(6, 6, 1) // two hits kill it; then a 1: back on top
	tb.Do(attackIntent(tb), "ragling")
	// On top of the monster deck, then drawn back into the empty slot.
	if top(tb, 0) != "ragling" {
		t.Error("Ragling did not come back from the top of the monster deck")
	}
	if n := len(tb.G.Players[0].Hand); n != 1 {
		t.Errorf("hand = %d, want 1: rewards come first", n)
	}
}
