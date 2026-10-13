package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheLeftHand(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"the_left_hand", "mystery_sack"}, Deactivated: []engine.CardRef{"mystery_sack"}},
		seat("cain"),
	}})
	tb.G.DealDamageTo(engine.Target{Player: 1, IsPlayer: true}, 2, 0, 0)
	tb.Settle("mystery_sack")
	if !tb.G.Object(tb.Find(0, "mystery_sack")).Charged {
		t.Error("Cain died: the sack should be recharged")
	}
}
