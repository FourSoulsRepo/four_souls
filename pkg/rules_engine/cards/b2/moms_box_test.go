package b2

import (
	"slices"
	"testing"
)

func TestMomsBox(t *testing.T) {
	tb := itemTable(t, items("moms_box", "mystery_sack"), "a_dime")
	tb.G.ForceRolls(4)
	tb.Activate(0, "mystery_sack", 0, "yes", "a_dime")
	if h := handCards(tb.G, 0); len(h) != 1 || h[0] == "a_dime" {
		t.Errorf("hand %v", h)
	}
}

// The answers are words, never card refs: "no" was the ref of No!, so
// the client showed that card on the decline button (19.42).
func TestMomsBoxAnswersAreNotCards(t *testing.T) {
	tb := itemTable(t, items("moms_box", "mystery_sack"), "a_dime")
	tb.G.ForceRolls(4)
	tb.Start(0, "mystery_sack", 0)
	if _, err := tb.G.Settle(); err == nil {
		t.Fatal("no question: the roll should ask")
	}
	got := tb.Options()
	if want := []string{"yes", "No"}; !slices.Equal(got, want) {
		t.Fatalf("options = %v, want %v", got, want)
	}
	for _, d := range Set.Cards {
		if slices.Contains(got, string(d.Ref)) {
			t.Errorf("option %q is the ref of a card", d.Ref)
		}
	}
}
