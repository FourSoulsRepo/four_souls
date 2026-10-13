package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTelepathyForDummies(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("cain"), {Character: "isaac", Items: items("telepathy_for_dummies")}},
		Monsters: items("trite", "fatty"),
	})
	tb.G.ForceRolls(6) // the start-of-turn roll on Isaac's turn
	tb.EndTurn()
	tb.G.ForceRolls(1, 6)
	tb.Do(attackIntent(tb), "trite", "yes") // the 1 becomes the stored 6
	if top(tb, 0) == "trite" {
		t.Error("the attack roll should become a 6 and kill the trite")
	}
	if d := tb.G.Players[1].Damage; d != 0 {
		t.Errorf("damage = %d, want 0", d)
	}
}
