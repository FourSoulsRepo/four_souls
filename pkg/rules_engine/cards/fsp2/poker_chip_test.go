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
