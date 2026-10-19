package r

import "testing"

func TestCharmedBony(t *testing.T) {
	// The roll of 6 first deals 1 damage to a player: Samson.
	tb := charmedGift(t, "charmed_bony", "player 3 (samson)")
	if tb.G.Players[1].Cents != 5 || tb.G.Players[2].Damage != 1 {
		t.Error("Samson takes 1 damage; Cain gains 5¢")
	}
}
