package b2

import "testing"

func TestXIIIDeath(t *testing.T) {
	tb := lootTable(t, "xiii_death")
	tb.Play(0, "xiii_death", foe)
	if !tb.G.Players[1].Dead {
		t.Error("the player is not dead")
	}
}
