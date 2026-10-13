package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCursedGlobin(t *testing.T) {
	tb := armed(t, items("mystery_sack"), "cursed_globin", "gurdy")
	tb.G.DealDamageTo(engine.Target{Object: topID(tb, 1)}, 3, 0, 0)
	tb.Settle()
	tb.G.ForceRolls(3)
	tb.Activate(0, "mystery_sack", 0)
	if d := tb.G.Object(topID(tb, 1)).Damage; d != 1 {
		t.Errorf("Gurdy damage = %d, want 3 - 2", d)
	}
}
