package fsp2

import (
	"slices"
	"testing"
)

func TestHourglass(t *testing.T) {
	tb := table(t, seat("isaac", "hourglass", "mystery_sack"), seat("cain", "boomerang"))
	tb.G.ForceRolls(2)
	tb.Activate(0, "mystery_sack", 0, "boomerang")
	if tb.G.Object(tb.Find(1, "boomerang")).Charged {
		t.Error("on a 2 an item may be deactivated")
	}
}

// Passive items are never offered, Hourglass itself included (19.64).
func TestHourglassSkipsPassiveItems(t *testing.T) {
	tb := table(t, seat("isaac", "hourglass", "mystery_sack"), seat("cain", "boomerang", "dads_lost_coin"))
	tb.G.ForceRolls(2)
	tb.Start(0, "mystery_sack", 0)
	if _, err := tb.G.Settle(); err == nil {
		t.Fatal("no question: the 2 should ask")
	}
	got := tb.Options()
	if !slices.Contains(got, "boomerang") {
		t.Errorf("options = %v: the charged Boomerang is offered", got)
	}
	for _, card := range []string{"hourglass", "dads_lost_coin"} {
		if slices.Contains(got, card) {
			t.Errorf("options = %v: passive %s is offered", got, card)
		}
	}
}
