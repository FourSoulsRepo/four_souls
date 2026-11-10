package challenge_test

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// R-KRAMPUS-08 (Ultra): a present at 0 brings the top 4 monster cards as
// minions.
func TestWhatsInTheBoxUltra(t *testing.T) {
	tb := krampus(t, engine.DifficultyUltra, engine.SituationPlayer{}, engine.SituationPlayer{})
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter", "dip", "squirt", "fly")
	tb.G.AddPresentCounters(engine.PresentTreasure, -3, engine.NoPlayer)
	tb.G.Run()
	tb.Settle()
	if len(tb.G.Minions) != 4 || !tb.G.Challenge.Presents[2].Destroyed {
		t.Errorf("minions %v", minionCards(tb.G))
	}
	check(t, tb.G)
}

// R-KRAMPUS-10, Q-30 (Ultra): Happy Holidays! comes at 7 souls, not before.
func TestHappyHolidaysUltra(t *testing.T) {
	tb := krampus(t, engine.DifficultyUltra, engine.SituationPlayer{Souls: souls(7 - 2)}, engine.SituationPlayer{})
	boost(tb, 0, 5)
	tb.G.ForceRolls(6)
	tb.RevealFromDeck("mixed_bag", "ask.no", "ask.no", "ask.no")
	if tb.G.Challenge.Holidays {
		t.Fatalf("Happy Holidays! at %d souls", tb.G.SoulsInPlay())
	}
	tb.G.Turn.Attacks = 1
	tb.G.ForceRolls(6)
	cain := player(tb, 1)
	tb.RevealFromDeck("mixed_bag", "ask.no", "ask.no", "ask.no", cain, cain, cain, cain)
	if !tb.G.Challenge.Holidays || len(tb.G.Players[1].InPlay) != 4 {
		t.Errorf("at %d souls: holidays %v, Cain's play area %d", tb.G.SoulsInPlay(), tb.G.Challenge.Holidays, len(tb.G.Players[1].InPlay))
	}
	check(t, tb.G)
}
