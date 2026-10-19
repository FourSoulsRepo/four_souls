package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHypercoagulation(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("the_dauntless", "hypercoagulation")}})
	h := tb.Find(1, "hypercoagulation")
	tb.G.AddCountersTo(h, 2)
	tb.EndTurn()
	if tb.G.PlayerHP(1) != 4 {
		t.Errorf("3 counters: HP %d, want 1 + 3", tb.G.PlayerHP(1))
	}
	tb.EndTurn()
	if tb.G.Object(h).CountersOf("") != 0 || len(tb.G.Players[1].Hand) != 1+3 {
		t.Error("at the end of his turn: counters off, loot 3")
	}
}
