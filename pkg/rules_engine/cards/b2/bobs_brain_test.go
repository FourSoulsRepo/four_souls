package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBobsBrain(t *testing.T) {
	tb := itemTable(t, items("bobs_brain"))
	tb.G.ForceRolls(3, 6) // the brain's roll: 3, deal 1 to a player; then the attack roll
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe, "fly")
	if d := tb.G.Players[1].Damage; d != 1 {
		t.Errorf("Cain's damage = %d, want 1", d)
	}
}
