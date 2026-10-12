package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCrookedPenny(t *testing.T) {
	for _, c := range []struct{ roll, want int }{{2, 10}, {5, 0}} {
		tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
			{Character: "isaac", Items: []engine.CardRef{"crooked_penny"}},
			{Character: "cain", Cents: 5},
		}})
		tb.G.ForceRolls(c.roll)
		tb.Activate(0, "crooked_penny", 0, foe)
		if got := tb.G.Players[1].Cents; got != c.want {
			t.Errorf("roll %d: Cain's cents = %d, want %d", c.roll, got, c.want)
		}
	}
}
