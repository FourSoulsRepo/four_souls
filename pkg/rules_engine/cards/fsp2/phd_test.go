package fsp2

import (
	"slices"
	"testing"
)

func TestPHD(t *testing.T) {
	tb := table(t, seat("isaac", "phd", "mystery_sack"), seat("cain"))
	tb.G.ForceRolls(2)
	tb.Activate(0, "mystery_sack", 0, "+1") // a 2 becomes a 3, once (R-DICE-10)
	if c := tb.G.Players[0].Cents; c != 4 {
		t.Errorf("cents = %d: the sack's 3 gives 4¢", c)
	}
}

// The answers are words, never card refs: "no" was the ref of No!, so
// the client showed that card on the decline button (19.42).
func TestPHDAnswersAreNotCards(t *testing.T) {
	tb := table(t, seat("isaac", "phd", "mystery_sack"), seat("cain"))
	tb.G.ForceRolls(2)
	tb.Start(0, "mystery_sack", 0)
	if _, err := tb.G.Settle(); err == nil {
		t.Fatal("no question: the roll should ask")
	}
	got := tb.Options()
	if want := []string{"+1", "-1", "No"}; !slices.Equal(got, want) {
		t.Fatalf("options = %v, want %v", got, want)
	}
	for _, s := range sets {
		for _, d := range s.Cards {
			if slices.Contains(got, string(d.Ref)) {
				t.Errorf("option %q is the ref of a card", d.Ref)
			}
		}
	}
}
