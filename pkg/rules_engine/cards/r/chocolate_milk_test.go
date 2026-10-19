package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestChocolateMilk(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("isaac", "chocolate_milk")}})
	milk := tb.Find(1, "chocolate_milk")
	tb.G.AddCountersTo(milk, 1)
	tb.EndTurn()
	if tb.G.PlayerATK(1) != 2 {
		t.Errorf("2 counters: ATK %d, want 2", tb.G.PlayerATK(1))
	}
	tb.EndTurn()
	tb.EndTurn("player 1 (cain)") // the 3rd counter: all removed, 1 damage
	if got := tb.G.Object(milk).CountersOf(""); got != 0 {
		t.Errorf("counters = %d, want 0", got)
	}
}
