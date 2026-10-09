package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheHaunt(t *testing.T) {
	tb := slayTable(t, "the_haunt", seat("isaac"), seat("cain"))
	// The second hit gives it +1 DC (5): the 4 then misses.
	tb.G.ForceRolls(6, 6, 4, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "the_haunt")
	if d := tb.G.Players[0].Damage; d != 1 {
		t.Errorf("damage = %d, want 1: the 4 missed", d)
	}
}
