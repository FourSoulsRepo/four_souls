package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The HP back shows as a heal of 1 (19.2). The prevented death is in
// testdata/situations/guppys_collar/.
func TestGuppysCollarHealsBack1HP(t *testing.T) {
	tb := itemTable(t, items("guppys_collar"), "xiii_death")
	tb.G.ForceRolls(2)
	ev := tb.Play(0, "xiii_death", me)
	healed := 0
	for _, e := range ev {
		if e.Kind == engine.EvHealed && e.Player == 0 {
			healed += e.Amount
		}
	}
	if healed != 1 {
		t.Errorf("healed %d, want 1 (R-MECH-46: to at least 1 HP)", healed)
	}
}
