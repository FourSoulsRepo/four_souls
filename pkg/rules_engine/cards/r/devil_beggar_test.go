package r

import "testing"

func TestDevilBeggar(t *testing.T) {
	tb := roomTable(t, "devil_beggar")
	n := len(tb.G.Players[0].Hand)
	tb.G.ForceRolls(1)
	useRoom(tb, 0, 0)
	if got := len(tb.G.Players[0].Hand) - n; got != 2 || tb.G.Players[0].Damage != 1 {
		t.Errorf("looted %d, damage %d; want 2 and 1", got, tb.G.Players[0].Damage)
	}
}
