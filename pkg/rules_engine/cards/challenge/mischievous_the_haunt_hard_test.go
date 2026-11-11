package challenge_test

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// R-TRICK-03, R-TRICK-06, R-TRICK-07 (Hard): 14 trick counters guard
// the boss; damage puts a fear counter on the damaged player only.
func TestHauntHard(t *testing.T) {
	tb := trick(t, engine.DifficultyHard, engine.SituationPlayer{Souls: souls(13)}, engine.SituationPlayer{})
	if tricks(tb) != 14 {
		t.Fatalf("trick counters %d, want 14", tricks(tb))
	}
	if opts := attackOptions(tb); slices.Contains(opts, "mischievous_the_haunt_hard") {
		t.Errorf("14 tricks, 13 souls: attack options %v", opts)
	}
	tb = trick(t, engine.DifficultyHard, engine.SituationPlayer{}, engine.SituationPlayer{})
	tb.G.DealDamageTo(engine.Target{Player: 1, IsPlayer: true}, 2, 0, 0)
	tb.G.Run()
	settle(tb)
	if fear(tb, 0) != 0 || fear(tb, 1) != 1 {
		t.Errorf("fear %d and %d, want 0 and 1", fear(tb, 0), fear(tb, 1))
	}
	check(t, tb.G)
}
