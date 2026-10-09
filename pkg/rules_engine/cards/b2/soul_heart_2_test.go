package b2

import "testing"

func TestSoulHeart(t *testing.T) {
	tb := lootTable(t, "soul_heart_2")
	tb.Play(0, "soul_heart_2", me)
	tb.Attack("fly", 1, 1, 6) // the first 1 damage is prevented, the second is not
	if got := tb.G.Players[0].Damage; got != 1 {
		t.Errorf("damage = %d, want 1", got)
	}
}
