package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGetOutOfJailCard(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"get_out_of_jail_card"}},
		{Character: "cain", Items: []engine.CardRef{"mystery_sack"}},
	}})
	tb.Play(0, "get_out_of_jail_card")
	tb.Start(0, "isaac", 0)
	tb.Pass(0)
	if rule := tb.Refused(1, "mystery_sack", 0); rule == "" {
		t.Error("Cain can't activate items this turn")
	}
	tb.Settle()
	tb.EndTurn()
	tb.Activate(1, "mystery_sack", 0)
}
