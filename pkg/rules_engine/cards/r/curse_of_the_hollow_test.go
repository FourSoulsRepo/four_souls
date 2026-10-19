package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCurseOfTheHollow(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"bomb", "a_nickel"}, Cents: 5, Souls: []engine.CardRef{"monstro", "gurdy"}}, seat("cain"),
	}})
	giveCurse(tb, 0, "curse_of_the_hollow")
	tb.EndTurn("bomb", "a_nickel")
	if len(tb.G.Players[0].Hand) != 0 || tb.G.Players[0].Cents != 3 {
		t.Error("2 souls: discard 2 and lose 2¢")
	}
}
