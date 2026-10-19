package r

import "testing"

func TestCharmedFatBat(t *testing.T) {
	if tb := charmedGift(t, "charmed_fat_bat"); len(tb.G.Players[1].InPlay) != 1 {
		t.Error("Cain gains +1 treasure")
	}
}
