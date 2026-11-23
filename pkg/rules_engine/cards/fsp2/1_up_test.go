package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Test1UpHealsOnce: one popup for the whole heal (19.2); the healed event's
// amount is an event field no situation check reads. The prevented death
// is in testdata/situations/1_up/.
func Test1UpHealsOnce(t *testing.T) {
	tb := armed(t, items("1_up"), "leech", "fatty")
	tb.G.ForceRolls(1) // the leech hits Isaac for 2: lethal
	ev := tb.Do(attackIntent(tb), "leech")
	// Isaac had 2 HP: heal to full is 2.
	var heals []int
	for _, e := range ev {
		if e.Kind == engine.EvHealed && e.Player == 0 {
			heals = append(heals, e.Amount)
		}
	}
	if len(heals) != 1 || heals[0] != 2 {
		t.Errorf("heals = %v, want one of 2", heals)
	}
}
