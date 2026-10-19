package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBloodLustR(t *testing.T) {
	tb := roomTable(t, "blood_lust_r")
	tb.EndTurn()
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 1, Kind: engine.IntentAttack}, "fly")
	if tb.G.Turn.MustAttacks != 1 || tb.G.Turn.Attacks != 1 {
		t.Errorf("after the first attack Cain must attack again: must %d, attacks %d", tb.G.Turn.MustAttacks, tb.G.Turn.Attacks)
	}
}
