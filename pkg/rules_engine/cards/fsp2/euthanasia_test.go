package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestEuthanasia(t *testing.T) {
	tb := armed(t, items("euthanasia"), "gurdy", "fatty")
	tb.G.ForceRolls(4, 4) // two 4s in a row: Gurdy (HP 5) dies
	ev := tb.Do(attackIntent(tb), "gurdy")
	if top(tb, 0) == "gurdy" {
		t.Error("Gurdy should die to the second 4")
	}
	if n := attackRolls(ev); n != 2 {
		t.Errorf("attack rolls = %d, want 2", n)
	}
}

// R-DICE-09: another roll between two attack rolls does not break them.
func TestEuthanasiaOtherRollsBetween(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: items("euthanasia", "mystery_sack")}, seat("cain")},
		Monsters: items("gurdy", "fatty"),
	})
	tb.G.ForceRolls(4, 2) // attack 4 hits; then the sack rolls 2
	declareAttack(tb, "gurdy")
	tb.Activate(0, "mystery_sack", 0)
	tb.G.ForceRolls(4)
	tb.Settle()
	if top(tb, 0) == "gurdy" {
		t.Error("two attack rolls of 4 in a row: Gurdy dies")
	}
}
