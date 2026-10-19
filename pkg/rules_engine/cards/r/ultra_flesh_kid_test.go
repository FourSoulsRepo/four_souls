package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestUltraFleshKid(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "ultra_flesh_kid"), seat("cain")},
		Monsters: []engine.CardRef{"fly"},
	})
	kid := tb.Find(0, "ultra_flesh_kid")
	tb.G.AddCountersTo(kid, 10)
	tb.Attack("fly", 6)
	if got := tb.G.Object(kid).CountersOf(""); got != 11 {
		t.Fatalf("counters = %d, want 11", got)
	}
	tb.EndTurn()
	if tb.G.SoulValue(0) != 0 {
		t.Error("11 counters: not yet a soul")
	}
	tb.G.AddCountersTo(kid, 1)
	tb.EndTurn()
	tb.EndTurn()
	if tb.G.SoulValue(0) != 1 {
		t.Error("12 counters at the end of his turn: a soul")
	}
}
