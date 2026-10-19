package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGoldenTrinket(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"jawbone"}, Hand: []engine.CardRef{"golden_trinket"}}, {Character: "cain", Cents: 5}},
	})
	tb.Play(0, "golden_trinket", "jawbone")
	if tb.G.Object(tb.Find(0, "jawbone")).CountersOf(gold) != 1 {
		t.Fatal("the Jawbone has no gold counter")
	}
	tb.Activate(0, "golden_trinket", 1, foe) // the Jawbone's ↷
	if tb.G.Players[1].Cents != 2 {
		t.Errorf("Cain's cents = %d, want 2: the trinket has the Jawbone's ↷", tb.G.Players[1].Cents)
	}
}
