package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPokerChip(t *testing.T) {
	for _, c := range []struct{ roll, want int }{{2, 1}, {5, 6}} {
		plays := 2
		tb := setup(t, engine.SituationSetup{
			Players:   []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"poker_chip"}, Hand: []engine.CardRef{"three_cents", "three_cents"}}, seat("cain")},
			LootPlays: &plays,
		})
		tb.G.ForceRolls(c.roll)
		tb.Play(0, "three_cents")
		if got := tb.G.Players[0].Cents; got != c.want {
			t.Errorf("roll %d: cents = %d, want %d", c.roll, got, c.want)
		}
		tb.Play(0, "three_cents") // only the first gain each turn
		if got := tb.G.Players[0].Cents; got != c.want+3 {
			t.Errorf("roll %d: second gain = %d, want +3", c.roll, got-c.want)
		}
	}
}

// Diplopia copying Poker Chip rolls with Poker Chip's ability (a
// simulation found this).
func TestPokerChipThroughDiplopia(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("diplopia"), Hand: items("three_cents")},
		{Character: "cain", Items: items("poker_chip")},
	}})
	tb.Activate(0, "diplopia", 0, "poker_chip")
	tb.G.ForceRolls(5)
	tb.Play(0, "three_cents")
	if c := tb.G.Players[0].Cents; c != 6 {
		t.Errorf("cents = %d, want double 3", c)
	}
}
