package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPlaydoughCookie(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("isaac", "playdough_cookie", "spoon_bender")}})
	tb.EndTurn("none")
	if tb.G.Players[1].Cents != 1 {
		t.Error("no item chosen: gain 1¢")
	}
	tb.EndTurn()
	tb.EndTurn("spoon_bender", "add a counter")
	if tb.G.Object(tb.Find(1, "spoon_bender")).CountersOf("") != 1 {
		t.Error("a counter on the Spoon Bender")
	}
}
