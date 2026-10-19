package r

import "testing"

func TestShellGame(t *testing.T) {
	tb := roomTable(t, "shell_game")
	tb.G.Players[0].Cents = 3
	tb.G.ForceRolls(5)
	useRoom(tb, 0, 0)
	if got := tb.G.Players[0].Cents; got != 5 {
		t.Errorf("cents = %d, want 5", got)
	}
}
