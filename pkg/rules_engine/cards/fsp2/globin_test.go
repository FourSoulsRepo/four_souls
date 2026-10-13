package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGlobin(t *testing.T) {
	tb := fight(t, "globin", "fly")
	tb.G.ForceRolls(5) // at 1 HP: heal 2
	tb.G.DealDamageTo(engine.Target{Object: topID(tb, 0)}, 3, 0, 0)
	tb.Settle()
	if d := tb.G.Object(topID(tb, 0)).Damage; d != 1 {
		t.Errorf("damage = %d, want 3 - 2", d)
	}
}
