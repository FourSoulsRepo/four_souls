package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestThePoop(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac", "the_poop"), seat("cain"))
	tb.G.ForceRolls(1, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly")
	if c := tb.G.Object(tb.Find(0, "the_poop")).CountersOf(""); c != 1 {
		t.Fatalf("counters = %d, want 1", c)
	}
	tb.Activate(0, "the_poop", 1)
	if len(tb.G.Shields) != 1 {
		t.Error("no shield")
	}
}
