package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// TestGodheadAsksForDiceValues: Godhead's two options are dice values
// (ask.dice), so the client draws them as dice (18.27). That they are 1
// and 6 is in godhead/offers-1-or-6.json.
func TestGodheadAsksForDiceValues(t *testing.T) {
	tb := ownRoll(t, 3, seat("isaac", "godhead"), seat("cain"))
	tb.Start(0, "godhead", 0)
	tb.Choose("ask.roll_of 3")
	tb.UntilChoose()
	opts := tb.G.Waiting.Options
	if len(opts) == 0 {
		t.Fatal("no options")
	}
	for i, o := range opts {
		if o.Key != engine.AskDice {
			t.Errorf("option %d key = %s, want ask.dice", i+1, o.Key)
		}
	}
}

// TestDieFacesAreDiceValues: "choose a number" before a roll (Crystal
// Ball, I - The Magician, Red Candle) offers the six dice values (18.27).
func TestDieFacesAreDiceValues(t *testing.T) {
	for i, o := range cardkit.DieFaces(nil, nil) {
		if o.Key != engine.AskDice {
			t.Errorf("face %d key = %s, want ask.dice", i+1, o.Key)
		}
	}
}
