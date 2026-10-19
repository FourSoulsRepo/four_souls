package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSpoilsOfWar(t *testing.T) {
	tb := roomTable(t, "spoils_of_war")
	n := len(tb.G.Players[1].Hand)
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 9, 1, 0)
	settleAny(tb)
	if got := len(tb.G.Players[1].Hand) - n; got != 1 {
		t.Errorf("Cain looted %d, want 1", got)
	}
}
