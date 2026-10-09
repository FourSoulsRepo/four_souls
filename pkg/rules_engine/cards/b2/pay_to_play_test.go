package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestPayToPlay(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("pay_to_play"), Cents: 10}, seat("cain", "breakfast"),
	}}, Set)
	tb.Activate(0, "pay_to_play", 0, "breakfast")
	if !hasItem(tb.G, 0, "breakfast") || tb.G.Players[0].Cents != 0 {
		t.Error("not stolen for 10¢")
	}
}
