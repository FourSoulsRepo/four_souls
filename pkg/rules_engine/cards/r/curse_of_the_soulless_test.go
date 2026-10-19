package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCurseOfTheSoulless(t *testing.T) {
	tb := fight(t, "monstro")
	giveCurse(tb, 0, "curse_of_the_soulless")
	hurt(tb, 0, 1)
	tb.Attack("monstro", 6)
	if tb.G.SoulValue(0) != 0 || !inDiscard(tb.G, engine.MonsterDeck, "monstro") {
		t.Error("he can't gain souls: Monstro goes to discard")
	}
}

func TestCurseOfTheSoullessPassedOn(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"mega_bomb"}}, seat("isaac"), seat("samson"),
	}})
	giveCurse(tb, 1, "curse_of_the_soulless")
	tb.Play(0, "mega_bomb", "player 2 (isaac)", "player 3 (samson)")
	if len(tb.G.Players[2].InPlay) != 1 {
		t.Error("when Isaac dies the curse goes to Samson, not to discard")
	}
}
