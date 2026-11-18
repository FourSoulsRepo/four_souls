package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// TestGodheadAsksForDiceValues: "change the result of a dice roll to a 1
// or 6" offers two dice values, 1 and 6, so the client draws them as
// dice (card godhead, 18.27).
func TestGodheadAsksForDiceValues(t *testing.T) {
	tb := ownRoll(t, 3, seat("isaac", "godhead"), seat("cain"))
	tb.Start(0, "godhead", 0)
	tb.Choose("ask.roll_of 3")
	tb.UntilChoose()
	opts := tb.G.Waiting.Options
	if len(opts) != 2 {
		t.Fatalf("options = %v, want 1 and 6", tb.Options())
	}
	for i, want := range []string{"1", "6"} {
		if opts[i].Key != engine.AskDice || tb.Options()[i] != want {
			t.Errorf("option %d = %s %q, want ask.dice %q", i, opts[i].Key, tb.Options()[i], want)
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
