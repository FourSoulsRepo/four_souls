package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGuppysCollar(t *testing.T) {
	tb := itemTable(t, items("guppys_collar"), "xiii_death")
	tb.G.ForceRolls(2)
	ev := tb.Play(0, "xiii_death", me)
	if tb.G.Players[0].Dead {
		t.Error("a 2 should prevent the death")
	}
	// The HP back shows as a heal (19.2).
	healed := 0
	for _, e := range ev {
		if e.Kind == engine.EvHealed && e.Player == 0 {
			healed += e.Amount
		}
	}
	if want := tb.G.PlayerHP(0); healed != 1 || want != 1 {
		t.Errorf("healed %d to %d HP, want 1 to 1", healed, want)
	}
}
