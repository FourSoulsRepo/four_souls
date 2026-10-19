package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBloodMoney(t *testing.T) {
	tb := roomTable(t, "blood_money")
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 9, 1, 0)
	settleAny(tb)
	if got := tb.G.Players[1].Cents; got != 4 {
		t.Errorf("Cain's cents = %d, want 4", got)
	}
}
