package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestTelepathyForDummiesShowsTheRoll: the view shows the stored roll
// until it is used (15.67). The changed attack roll is in
// testdata/situations/telepathy_for_dummies/.
func TestTelepathyForDummiesShowsTheRoll(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("cain"), {Character: "isaac", Items: items("telepathy_for_dummies")}},
		Monsters: items("trite", "fatty"),
	})
	tb.G.ForceRolls(6) // the start-of-turn roll on Isaac's turn
	tb.EndTurn()
	if d := viewDice(tb, 1, "telepathy_for_dummies"); len(d) != 1 || d[0] != 6 {
		t.Errorf("view dice = %v, want the stored [6]", d)
	}
	tb.G.ForceRolls(1, 6)
	tb.Do(attackIntent(tb), "trite", "ask.yes") // the 1 becomes the stored 6
	if d := viewDice(tb, 1, "telepathy_for_dummies"); d != nil {
		t.Errorf("view dice = %v after the roll, want none", d)
	}
}
