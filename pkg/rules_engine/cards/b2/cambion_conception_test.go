package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCambionConception(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac", "cambion_conception"), seat("cain"))
	tb.G.Object(tb.Find(0, "cambion_conception")).Counters = []engine.Counter{{Count: 5}}
	tb.G.ForceRolls(1, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly") // 1 damage: 6 counters
	if n := len(tb.G.Players[0].InPlay); n != 2 {
		t.Errorf("%d items, want the conception and a treasure", n)
	}
	if c := tb.G.Object(tb.Find(0, "cambion_conception")).CountersOf(""); c != 0 {
		t.Errorf("counters = %d, want 0", c)
	}
}
