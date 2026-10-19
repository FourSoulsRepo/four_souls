package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestStaticShock(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "sack_of_pennies"), seat("cain")},
		Rooms:   []engine.CardRef{"static_shock"},
	})
	tb.Activate(0, "sack_of_pennies", 0)
	settleAny(tb)
	if got := tb.G.Players[0].Damage; got != 1 {
		t.Errorf("Isaac's damage = %d, want 1", got)
	}
}
