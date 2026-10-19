package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSoulOfWrath(t *testing.T) {
	tb := bonusTable(t, "soul_of_wrath")
	tb.G.AddCountersTo(tb.G.BonusSouls[0], 5)
	tb.EndTurn()                                                         // Cain's turn
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 9, 1, 0) // Isaac dies
	settleAny(tb)
	if s := tb.G.SoulValue(1); s != 1 {
		t.Errorf("Cain's soul value = %d, want 1: the active player gains it on the 6th death", s)
	}
}
