package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSacrificialDagger(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "sacrificial_dagger"), seat("cain")},
		Monsters: []engine.CardRef{"pooter", "fly"},
	})
	tb.Attack("pooter", 1, 6, 6)
	dagger := tb.Find(0, "sacrificial_dagger")
	if got := tb.G.Object(dagger).CountersOf(""); got != 1 {
		t.Fatalf("counters = %d, want 1", got)
	}
	tb.G.AddCountersTo(dagger, 2)
	tb.Activate(0, "sacrificial_dagger", 1, foe)
	if tb.G.Players[1].Damage != 1 {
		t.Error("1 damage to Cain")
	}
}
