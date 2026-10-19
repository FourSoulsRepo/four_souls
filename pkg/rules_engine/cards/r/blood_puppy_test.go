package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBloodPuppy(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"blood_puppy"}, Hand: []engine.CardRef{"mega_bomb"}}, seat("cain"),
	}})
	puppy := tb.Find(0, "blood_puppy")
	tb.G.AddNamedCounters(puppy, level, 3)
	tb.Play(0, "mega_bomb", foe)
	if tb.G.PlayerATK(0) != 2 {
		t.Errorf("ATK = %d, want 2 at level 4", tb.G.PlayerATK(0))
	}
}
