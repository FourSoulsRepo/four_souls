package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestMomsCoinPurse(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac"), seat("cain", "moms_coin_purse")}}, Set)
	tb.EndTurn()
	if n := len(tb.G.Players[1].Hand); n != 2 {
		t.Errorf("looted %d in the loot step, want 2", n)
	}
}
