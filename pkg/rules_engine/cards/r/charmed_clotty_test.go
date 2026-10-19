package r

import "testing"

func TestCharmedClotty(t *testing.T) {
	if tb := charmedGift(t, "charmed_clotty"); len(tb.G.Players[1].Hand) != 1 {
		t.Error("Cain loots 1")
	}
}
