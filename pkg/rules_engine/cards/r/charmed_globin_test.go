package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCharmedGlobin(t *testing.T) {
	tb := fight(t, "charmed_globin")
	hurt(tb, 0, 1)
	tb.G.ForceRolls(1, 6, 6) // a miss heals it to 2 HP: two hits needed
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "charmed_globin", foe)
	if tb.G.Turn.AttackRolls != 3 {
		t.Errorf("attack rolls = %d, want 3", tb.G.Turn.AttackRolls)
	}
}
