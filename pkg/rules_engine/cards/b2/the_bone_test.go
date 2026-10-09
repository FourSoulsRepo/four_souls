package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheBone(t *testing.T) {
	tb := monsterTable(t, []engine.CardRef{"fly"}, seat("the_forgotten", "the_bone"), seat("cain"))
	bone := tb.Find(0, "the_bone")
	if rule := tb.Refused(0, "the_bone", 2); rule != "R-ABIL-07" {
		t.Errorf("removing counters without any: refused by %s, want R-ABIL-07", rule)
	}
	tb.Activate(0, "the_bone", 0)
	if got := tb.G.Object(bone).CountersOf(""); got != 1 {
		t.Fatalf("counters = %d, want 1", got)
	}

	tb.G.Object(bone).Counters = []engine.Counter{{Count: 7}}
	fly, _ := tb.G.Monsters[0].TopOf()
	tb.Activate(0, "the_bone", 2, "fly")
	if tb.G.Object(fly).Zone.Kind == engine.ZoneInPlay {
		t.Error("the fly survived 1 damage")
	}
	if got := tb.G.Object(bone).CountersOf(""); got != 5 {
		t.Fatalf("counters = %d, want 5", got)
	}

	tb.Activate(0, "the_bone", 3)
	if got := tb.G.SoulValue(0); got != 1 {
		t.Errorf("soul value = %d, want 1", got)
	}
	if rule := tb.Refused(0, "the_bone", 0); rule != "R-ABIL-09" {
		t.Errorf("a soul's ability refused by %s, want R-ABIL-09 (no abilities)", rule)
	}
}

func TestTheBoneAddsToARoll(t *testing.T) {
	tb := rollTable(t, 3, seat("the_forgotten", "the_bone"), seat("cain"))
	tb.G.Object(tb.Find(0, "the_bone")).Counters = []engine.Counter{{Count: 1}}
	tb.Pass(1)
	ev := tb.Activate(0, "the_bone", 1, "roll of 3")
	if got := resolvedRoll(t, ev); got != 4 {
		t.Errorf("roll resolved as %d, want 4", got)
	}
}
