package r

import "testing"

func TestSocialGoals(t *testing.T) {
	tb := roomTable(t, "social_goals")
	id := room(tb)
	for _, goal := range []string{"1. Play 5 loot cards.", "2. Kill 3 monsters.", "4. Purchase 3 items."} {
		tb.G.AddNamedCounters(id, goal, 4)
	}
	tb.G.AddNamedCounters(id, goalsDone, 3)
	tb.G.Players[1].Cents = 6
	if tb.G.GiveCents(1, 0, 6) != 6 {
		t.Fatal("Cain gives Isaac 6¢")
	}
	tb.G.Run()
	settleAny(tb)
	if len(itemsOf(tb.G, 0)) != 2 || len(itemsOf(tb.G, 1)) != 2 {
		t.Error("the 4th goal: each player gains +2 treasure")
	}
}
