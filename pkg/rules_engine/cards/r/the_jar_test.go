package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheJar(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"the_jar"}, Hand: []engine.CardRef{"bomb"}, Cents: 4}, seat("cain"),
	}})
	tb.Activate(0, "the_jar", 0, foe)
	tb.Play(0, "bomb", foe)
	if tb.G.Players[1].Damage != 0 || tb.G.Players[0].Cents != 0 {
		t.Error("the bomb's 1 damage is prevented, for 4¢")
	}
}
