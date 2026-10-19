package r

import "testing"

func TestMightForTheMeek(t *testing.T) {
	tb := roomTable(t, "might_for_the_meek")
	tb.G.Players[1].InPlay = append(tb.G.Players[1].InPlay, soulOf(tb, 1, "chub"))
	if tb.G.PlayerATK(0) != 2 || tb.G.PlayerATK(1) != 1 {
		t.Errorf("ATK %d and %d, want 2 and 1", tb.G.PlayerATK(0), tb.G.PlayerATK(1))
	}
	tb.Attack("fly", 6)
	tb.Attack("leech", 6, 6)
	if tb.G.Turn.AttackRolls < 2 {
		t.Error("Isaac (fewest souls) may attack an additional time")
	}
}
