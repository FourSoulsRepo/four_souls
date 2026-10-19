package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPunchingBag(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac"), seat("cain", "punching_bag")}})
	if tb.G.PlayerHP(1) != 4 {
		t.Errorf("Cain's HP = %d, want 2 + 2", tb.G.PlayerHP(1))
	}
	tb.G.ForceRolls(5, 5)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "punching_bag")
	if hasItem(tb.G, 1, "punching_bag") {
		t.Fatal("killed, it is destroyed (R-MECH-05)")
	}
	if tb.G.Players[0].Cents != 5 {
		t.Errorf("the active player gains its 5¢; cents %d", tb.G.Players[0].Cents)
	}
}
