package challenge_test

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Pilfering Krampus's guard, rolls, random attack and minions are in
// testdata/situations/pilfering_krampus_* (24.27). What stays here: the
// view of the random attack (18.35).

// R-KRAMPUS-13, Q-31, 18.35: the view says the random attack is owed,
// for every seat: "players must attack a random minion once each turn
// if able"; once it is made, it is no longer owed.
func TestKrampusForcedAttackView(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{}, "pooter", "dip")
	if !tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).RandomAttack {
		t.Error("view: the random attack is not owed")
	}
	boost(tb, 0, 4)
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack})
	if tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).RandomAttack {
		t.Error("view: the random attack is still owed after it was made")
	}
	check(t, tb.G)
}
