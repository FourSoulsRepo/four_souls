package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestIBS(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "the_soiled", Items: []engine.CardRef{"ibs"}, Hand: []engine.CardRef{"bomb", "a_nickel", "dice_shard"}}, seat("cain"),
	}})
	tb.Activate(0, "ibs", 0, "Gain 3¢.", "bomb")
	if tb.G.Players[0].Cents != 3 {
		t.Error("discard a card: gain 3¢")
	}
}
