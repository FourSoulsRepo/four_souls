package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestTheBoneSoulHasNoAbilities: the soul lost all abilities, so its ↷
// is refused. The engine cites R-ABIL-09; 19.226 asks for the rule that
// says it, then this moves to a file. The rest of The Bone is situation
// files in the_bone/.
func TestTheBoneSoulHasNoAbilities(t *testing.T) {
	tb := monsterTable(t, []engine.CardRef{"fly"}, seat("the_forgotten", "the_bone"), seat("cain"))
	tb.Activate(0, "the_bone", 0)
	tb.G.Object(tb.Find(0, "the_bone")).Counters = []engine.Counter{{Count: 5}}
	tb.Activate(0, "the_bone", 3)
	if got := tb.G.SoulValue(0); got != 1 {
		t.Fatalf("soul value = %d, want 1", got)
	}
	if rule := tb.Refused(0, "the_bone", 0); rule != "R-ABIL-09" {
		t.Errorf("a soul's ability refused by %s, want R-ABIL-09 (no abilities)", rule)
	}
}
