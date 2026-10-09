package b2

import "testing"

func TestFannyPack(t *testing.T) {
	tb := itemTable(t, items("fanny_pack"))
	tb.Attack("fly", 1, 6)
	if h := len(tb.G.Players[0].Hand); h != 1 {
		t.Errorf("hand = %d, want 1", h)
	}
}
