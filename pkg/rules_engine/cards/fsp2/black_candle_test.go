package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBlackCandle(t *testing.T) {
	tb := armed(t, items("black_candle"), "fatty", "monstro")
	tb.G.PutOnDeck(engine.MonsterDeck, "fly", "curse_of_pain", "leech", "trite", "dip", "spider")
	tb.Activate(0, "black_candle", 0, foe)
	found := false
	for _, id := range tb.G.Players[1].InPlay {
		found = found || (tb.G.Object(id).Card == "curse_of_pain" && tb.G.Object(id).Role == engine.RoleCurse)
	}
	if !found {
		t.Error("Cain should have the Curse of Pain")
	}
	for _, id := range tb.G.DeckTop(engine.MonsterDeck, 5) {
		if c := tb.G.Object(id).Card; c == "spider" || c == "fly" {
			t.Errorf("%s should be at the bottom", c)
		}
	}
}
