package challenge_test

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// R-TRICK-03, R-TRICK-07, Q-50 (Ultra Hard): 14 trick counters; each
// instance of damage to a player puts a fear counter on each player.
func TestHauntUltraFear(t *testing.T) {
	tb := trick(t, engine.DifficultyUltra, engine.SituationPlayer{}, engine.SituationPlayer{})
	if tricks(tb) != 14 {
		t.Errorf("trick counters %d, want 14", tricks(tb))
	}
	tb.G.DealDamageTo(engine.Target{Player: 1, IsPlayer: true}, 2, 0, 0)
	tb.G.Run()
	settle(tb)
	if fear(tb, 0) != 1 || fear(tb, 1) != 1 {
		t.Errorf("fear %d and %d, want 1 each", fear(tb, 0), fear(tb, 1))
	}
	check(t, tb.G)
}
