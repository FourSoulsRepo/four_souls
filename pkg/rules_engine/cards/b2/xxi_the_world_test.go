package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestXXITheWorld(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"xxi_the_world"}},
		{Character: "cain", Hand: []engine.CardRef{"a_dime"}},
	}}, Set)
	ev := tb.Play(0, "xxi_the_world")
	if !enginetest.Has(ev, engine.EvLookedAt, 0) {
		t.Error("no look at the other hand")
	}
	if n := len(tb.G.Players[0].Hand); n != 2 {
		t.Errorf("hand %d, want 2", n)
	}
}
