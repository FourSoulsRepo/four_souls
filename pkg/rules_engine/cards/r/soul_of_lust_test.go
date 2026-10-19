package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// bonusTable: Isaac (active) and Cain, a fly and a leech, and the active
// bonus soul.
func bonusTable(t *testing.T, soul engine.CardRef) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{
		Players:    []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Monsters:   []engine.CardRef{"fly", "leech"},
		BonusSouls: []engine.CardRef{soul},
	})
}

// bonusCounters returns the counters on the only active bonus soul.
func bonusCounters(tb *enginetest.Table) int {
	return tb.G.Object(tb.G.BonusSouls[0]).CountersOf("")
}

func TestSoulOfLust(t *testing.T) {
	tb := bonusTable(t, "soul_of_lust")
	tb.G.AddCountersTo(tb.G.BonusSouls[0], 5)
	if n := bonusCounters(tb); n != 5 {
		t.Fatalf("counters = %d, want 5", n)
	}
	tb.Attack("fly", 6)
	settleAny(tb)
	if s := tb.G.SoulValue(0); s != 1 {
		t.Errorf("Isaac's soul value = %d, want 1 (the 6th kill gives the soul)", s)
	}
	if len(tb.G.BonusSouls) != 0 {
		t.Error("the soul is still outside the game")
	}
}
