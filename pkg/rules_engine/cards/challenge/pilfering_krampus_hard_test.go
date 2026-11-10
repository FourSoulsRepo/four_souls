package challenge_test

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// R-KRAMPUS-12 (Hard): the table, one roll without minions.
func TestKrampusRollHard(t *testing.T) {
	for roll, want := range map[int][4]int{
		1: {7, 4, 3, 3}, 2: {7, 4, 3, 3}, 3: {12, 2, 3, 3}, 4: {12, 2, 3, 3}, 5: {12, 4, 2, 3}, 6: {12, 4, 3, 2},
	} {
		tb := krampus(t, engine.DifficultyHard, engine.SituationPlayer{}, engine.SituationPlayer{})
		tb.G.ForceRolls(roll, 1, 1)
		tb.EndTurn()
		if got := counters(tb.G); got != want {
			t.Errorf("roll %d: counters %v, want %v", roll, got, want)
		}
	}
}

// R-KRAMPUS-12, Q-41 (Hard): one more roll for every 2 minions in play,
// counted once as the ability resolves. Hard has no forced attack
// (R-KRAMPUS-13).
func TestKrampusExtraRollsHard(t *testing.T) {
	for _, c := range []struct {
		minions []engine.CardRef
		want    [4]int
	}{
		{[]engine.CardRef{"pooter"}, [4]int{12, 4, 3, 2}},                         // 1 roll: 6
		{[]engine.CardRef{"pooter", "dip"}, [4]int{12, 4, 2, 2}},                  // 2 rolls: 6, 5
		{[]engine.CardRef{"pooter", "dip", "squirt"}, [4]int{12, 4, 2, 2}},        // 2 rolls
		{[]engine.CardRef{"pooter", "dip", "squirt", "fly"}, [4]int{12, 2, 2, 2}}, // 3 rolls: 6, 5, 4
	} {
		tb := krampus(t, engine.DifficultyHard, engine.SituationPlayer{}, engine.SituationPlayer{}, c.minions...)
		if blocked(tb) {
			t.Error("Hard: the turn may end with a minion in play")
		}
		tb.G.ForceRolls(6, 5, 4, 1, 1, 1)
		tb.EndTurn()
		if got := counters(tb.G); got != c.want {
			t.Errorf("%d minions: counters %v, want %v", len(c.minions), got, c.want)
		}
		check(t, tb.G)
	}
}

// R-KRAMPUS-11 (Hard): guarded below 6 souls.
func TestKrampusGuardedHard(t *testing.T) {
	for n, want := range map[int]bool{5: false, 6: true} {
		tb := krampus(t, engine.DifficultyHard, engine.SituationPlayer{Souls: souls(n)}, engine.SituationPlayer{})
		if got := slices.Contains(attackOptions(tb), "pilfering_krampus_hard"); got != want {
			t.Errorf("%d souls: attackable %v", n, got)
		}
	}
}
