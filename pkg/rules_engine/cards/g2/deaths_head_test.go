package g2

import "testing"

func TestDeathsHead(t *testing.T) {
	tb := fight(t, "deaths_head", "fly")
	if atk := tb.G.MonsterATK(topID(tb, 1)); atk != 2 {
		t.Errorf("fly ATK = %d, want 1 + 1", atk)
	}
	tb.G.ForceRolls(6)
	tb.Do(attackIntent(tb), "fly") // Death's Head can't be attacked: not offered
	if top(tb, 0) == "deaths_head" {
		t.Error("Death's Head should die when the fly dies")
	}
	if n := len(tb.G.Players[0].Hand); n != 1 {
		t.Errorf("hand = %d, want 1: Death's Head's reward (1 loot)", n)
	}
}
