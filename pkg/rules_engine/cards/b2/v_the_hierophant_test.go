package b2

import "testing"

func TestVTheHierophant(t *testing.T) {
	tb := lootTable(t, "v_the_hierophant")
	tb.Play(0, "v_the_hierophant", me)
	tb.Attack("leech", 1, 6) // the leech's 2 damage is prevented
	if got := tb.G.Players[0].Damage; got != 0 {
		t.Errorf("damage = %d, want 0", got)
	}
}
