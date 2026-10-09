package b2

import "testing"

func TestGuppysPaw(t *testing.T) {
	tb := monsterTable(t, items("leech"), seat("isaac", "guppys_paw"), seat("cain"))
	tb.Activate(0, "guppys_paw", 0, me)
	if hp := tb.G.PlayerHP(0); hp != 1 {
		t.Errorf("HP after paying = %d, want 1", hp)
	}
	tb.Attack("leech", 1, 6) // the leech's 2 damage is prevented
	if tb.G.Players[0].Dead {
		t.Error("the shield did not hold")
	}
}
