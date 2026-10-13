package fsp2

import "testing"

func TestRainbowBaby(t *testing.T) {
	tb := table(t, seat("isaac", "rainbow_baby"), seat("cain"))
	tb.Activate(0, "rainbow_baby", 0, "Each player gains 2¢.")
	if tb.G.Players[0].Cents != 2 || tb.G.Players[1].Cents != 2 {
		t.Error("each player gains 2¢")
	}
}
