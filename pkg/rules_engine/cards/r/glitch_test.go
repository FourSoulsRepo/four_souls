package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestGlitchStaysEternalAndGainsAtTurnEnd reads engine state no situation
// check says: Glitch stays eternal with a gained ability (S-YUG), and the
// ability gained as the turn ends. The steal is in
// testdata/situations/glitch/.
func TestGlitchStaysEternalAndGainsAtTurnEnd(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "cain", Cents: 5}, seat("the_capricious", "glitch")},
		Shop:    []engine.CardRef{"jawbone", "spoon_bender"},
	})
	tb.EndTurn("jawbone") // his turn starts: Glitch gains the Jawbone's ↷
	if !tb.G.Eternal(tb.Find(1, "glitch")) {
		t.Error("it stays eternal (S-YUG)")
	}
	tb.EndTurn("spoon_bender")
	if tb.G.Object(tb.Find(1, "glitch")).Gained != "spoon_bender" {
		t.Error("as his turn ends it gains the Spoon Bender's abilities till his next turn")
	}
}
