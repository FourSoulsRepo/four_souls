package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGimpy(t *testing.T) {
	tb := table(t, seat("whore_of_babylon", "gimpy"), seat("cain"))
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 1, 1, 0)
	tb.Settle("gain 2¢")
	if c := tb.G.Players[0].Cents; c != 2 {
		t.Errorf("cents = %d, want 2", c)
	}
	tb2 := table(t, seat("whore_of_babylon", "gimpy"), seat("cain"))
	tb2.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 1, 1, 0)
	tb2.Settle("+2 ATK till end of turn")
	if atk := tb2.G.PlayerATK(0); atk != 3 {
		t.Errorf("ATK = %d, want 1 + 2", atk)
	}
}
