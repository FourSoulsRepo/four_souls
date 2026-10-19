package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGlitch(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "cain", Cents: 5}, seat("the_capricious", "glitch")},
		Shop:    []engine.CardRef{"jawbone", "spoon_bender"},
	})
	tb.EndTurn("jawbone") // his turn starts: Glitch gains the Jawbone's ↷
	tb.Activate(1, "glitch", 2, "player 1 (cain)")
	if tb.G.Players[1].Cents != 3 || tb.G.Players[0].Cents != 2 {
		t.Errorf("the gained ↷ steals 3¢ from Cain: cents %d and %d", tb.G.Players[0].Cents, tb.G.Players[1].Cents)
	}
	if !tb.G.Eternal(tb.Find(1, "glitch")) {
		t.Error("it stays eternal (S-YUG)")
	}
	tb.EndTurn("spoon_bender")
	if tb.G.Object(tb.Find(1, "glitch")).Gained != "spoon_bender" {
		t.Error("as his turn ends it gains the Spoon Bender's abilities till his next turn")
	}
}
