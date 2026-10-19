package r

import "testing"

func TestCharmedGreedling(t *testing.T) {
	if tb := charmedGift(t, "charmed_greedling"); tb.G.Players[1].Cents != 10 {
		t.Error("Cain gains 10¢")
	}
}
