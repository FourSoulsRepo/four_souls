package b2

import "testing"

func TestVITheLovers(t *testing.T) {
	tb := lootTable(t, "vi_the_lovers")
	tb.Play(0, "vi_the_lovers", me)
	if got := tb.G.PlayerHP(0); got != 4 {
		t.Errorf("HP = %d, want 4", got)
	}
	tb.EndTurn()
	if got := tb.G.PlayerHP(0); got != 2 {
		t.Errorf("HP after the turn = %d, want 2", got)
	}
}
