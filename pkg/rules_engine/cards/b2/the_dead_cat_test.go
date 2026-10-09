package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestTheDeadCat(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Cents: 10}, seat("cain")},
		Shop:     items("the_dead_cat"),
		Monsters: items("leech"),
	}, Set)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "the_dead_cat")
	cat := tb.Find(0, "the_dead_cat")
	if c := tb.G.Object(cat).CountersOf(""); c != 9 {
		t.Fatalf("counters = %d, want 9", c)
	}
	tb.Attack("leech", 1, 6)
	if d, c := tb.G.Players[0].Damage, tb.G.Object(cat).CountersOf(""); d != 0 || c != 7 {
		t.Errorf("damage %d counters %d, want 0 and 7", d, c)
	}
}
