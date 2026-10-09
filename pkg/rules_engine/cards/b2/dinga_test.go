package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDinga(t *testing.T) {
	tb := slayTable(t, "dinga", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6, 6, 6, 3)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "dinga")
	if c := tb.G.Players[0].Cents; c != 6 {
		t.Errorf("cents = %d, want 6 (a reward roll of 3, doubled)", c)
	}
}
