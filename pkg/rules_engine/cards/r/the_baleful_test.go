package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheBalefulCharacter(t *testing.T) { testCharacter(t, "the_baleful", "soulbond") }

func TestTheBalefulATK(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("the_baleful"), seat("cain")}})
	if tb.G.PlayerATK(0) != 1 {
		t.Error("no souls: ATK 1")
	}
	tb.G.Players[0].InPlay = append(tb.G.Players[0].InPlay, soulOf(tb, 0, "monstro"))
	if tb.G.PlayerATK(0) != 2 {
		t.Error("1 soul: +1 ATK")
	}
	tb.G.Players[0].InPlay = append(tb.G.Players[0].InPlay, soulOf(tb, 0, "gurdy"), soulOf(tb, 0, "larry_jr"))
	if tb.G.PlayerATK(0) != 3 {
		t.Error("3 souls: +2 ATK instead")
	}
}
