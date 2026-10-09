package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// rollTable has a roll of n on the stack, owned by player 1.
func rollTable(t *testing.T, n int, seats ...engine.SituationPlayer) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{
		Players: seats,
		Stack:   []engine.SituationStack{{DiceRoll: n, Owner: 1}},
	}, Set)
}

// resolvedRoll is the result of the first roll that resolved.
func resolvedRoll(t *testing.T, events []engine.Event) int {
	t.Helper()
	for _, e := range events {
		if e.Kind == engine.EvRollResolved {
			return e.Amount
		}
	}
	t.Fatal("no roll resolved")
	return 0
}

func TestTheD6(t *testing.T) {
	tb := rollTable(t, 2, seat("isaac", "the_d6"), seat("cain"))
	tb.Pass(1)
	tb.G.ForceRolls(5)
	ev := tb.Activate(0, "the_d6", 0, "roll of 2")
	if got := resolvedRoll(t, ev); got != 5 {
		t.Errorf("rerolled roll resolved as %d, want 5", got)
	}
	d6 := tb.Find(0, "the_d6")
	if tb.G.Object(d6).Charged {
		t.Fatal("the D6 is charged after use")
	}
	tb.EndTurn()
	if !tb.G.Object(d6).Charged {
		t.Error("the D6 did not recharge at the end of its controller's turn")
	}
}

func TestTheD6NeedsARoll(t *testing.T) {
	tb := table(t, seat("isaac", "the_d6"), seat("cain"))
	if rule := tb.Refused(0, "the_d6", 0); rule != "R-ABIL-06" {
		t.Errorf("refused by %s, want R-ABIL-06 (no dice roll)", rule)
	}
}
