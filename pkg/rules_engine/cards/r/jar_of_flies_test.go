package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestJarOfFlies(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "jar_of_flies"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	if rule := tb.Refused(0, "jar_of_flies", 1); rule != "R-ABIL-07" {
		t.Errorf("no counters: refused by %s", rule)
	}
	tb.G.AddCountersTo(tb.Find(0, "jar_of_flies"), 3)
	tb.Activate(0, "jar_of_flies", 1, "1")
	if len(tb.G.Players[0].Hand) != 1 {
		t.Error("1 counter: loot 1")
	}
	tb.G.Object(tb.Find(0, "jar_of_flies")).Charged = true
	tb.Activate(0, "jar_of_flies", 1, "2")
	if len(tb.G.Players[0].Hand) != 2 {
		t.Error("2 counters: loot 1")
	}
}

func TestJarOfFliesThreeCounters(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "jar_of_flies"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.AddCountersTo(tb.Find(0, "jar_of_flies"), 3)
	tb.Activate(0, "jar_of_flies", 1, "3", "fly")
	if top(tb, 0) == "fly" || len(tb.G.Players[0].Hand) != 0 {
		t.Error("3 counters: 1 damage to a monster instead of the loot")
	}
}

func TestJarOfFliesCounts(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "jar_of_flies"), seat("cain")},
		Monsters: []engine.CardRef{"fly"},
	})
	tb.Attack("fly", 6)
	if got := tb.G.Object(tb.Find(0, "jar_of_flies")).CountersOf(""); got != 1 {
		t.Errorf("counters = %d, want 1", got)
	}
}
