package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCurseOfBloodLust(t *testing.T) {
	tb := fight(t, "fatty", "monstro")
	tb.RevealFromDeck("curse_of_blood_lust", me)
	tb.EndTurn()
	tb.EndTurn()
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err == nil {
		t.Error("Isaac must attack before ending his turn")
	}
}
