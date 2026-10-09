package b2

import "testing"

func TestChaosCard(t *testing.T) {
	tb := itemTable(t, items("chaos_card"))
	tb.Activate(0, "chaos_card", 0, "Kill a player or monster.", foe)
	if !tb.G.Players[1].Dead || hasItem(tb.G, 0, "chaos_card") {
		t.Error("Cain should be dead and the card gone")
	}
}
