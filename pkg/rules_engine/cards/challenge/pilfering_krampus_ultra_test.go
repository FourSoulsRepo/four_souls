package challenge_test

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// R-KRAMPUS-12 (Ultra Hard): the table, one roll without minions.
func TestKrampusRollUltra(t *testing.T) {
	for roll, want := range map[int][4]int{
		1: {8, 4, 3, 3}, 2: {16, 2, 3, 3}, 3: {16, 2, 3, 3}, 4: {16, 4, 2, 3}, 5: {16, 4, 2, 3}, 6: {16, 4, 3, 2},
	} {
		tb := krampus(t, engine.DifficultyUltra, engine.SituationPlayer{}, engine.SituationPlayer{})
		tb.G.ForceRolls(roll, 1, 1)
		tb.EndTurn()
		if got := counters(tb.G); got != want {
			t.Errorf("roll %d: counters %v, want %v", roll, got, want)
		}
	}
}

// R-KRAMPUS-12, Q-41 (Ultra Hard): one more roll for each minion.
// Minions that the rolls bring add none this turn.
func TestKrampusExtraRollsUltra(t *testing.T) {
	tb := krampus(t, engine.DifficultyUltra, engine.SituationPlayer{}, engine.SituationPlayer{}, "pooter", "dip")
	tb.G.ForceRolls(6, 6, 4, 1, 1, 1)
	tb.EndTurn()
	if got := counters(tb.G); got != [4]int{16, 4, 2, 1} {
		t.Errorf("2 minions, 3 rolls: counters %v", got)
	}

	tb = krampus(t, engine.DifficultyUltra, engine.SituationPlayer{}, engine.SituationPlayer{})
	setCounters(tb, engine.PresentSoul, 1)
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter", "dip", "squirt", "fly")
	tb.G.ForceRolls(6, 1, 1, 1)
	tb.EndTurn()
	if got := counters(tb.G); got != [4]int{16, 4, 3, 0} || len(tb.G.Minions) != 4 {
		t.Errorf("one roll brought 4 minions and no more rolls: counters %v, minions %v", got, minionCards(tb.G))
	}
	check(t, tb.G)
}

// R-KRAMPUS-11 (Ultra Hard): guarded below 7 souls.
func TestKrampusGuardedUltra(t *testing.T) {
	for n, want := range map[int]bool{6: false, 7: true} {
		tb := krampus(t, engine.DifficultyUltra, engine.SituationPlayer{Souls: souls(n)}, engine.SituationPlayer{})
		if got := slices.Contains(attackOptions(tb), "pilfering_krampus_ultra"); got != want {
			t.Errorf("%d souls: attackable %v", n, got)
		}
	}
}
