package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCursedLilHaunt(t *testing.T) {
	tb := fight(t, "cursed_lil_haunt")
	tb.G.ForceRolls(1, 6, 6, 6) // two dice each roll
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "1", "6")
	if tb.G.Players[0].Damage != 1 || top(tb, 1) == "fly" {
		t.Error("Cain picks the 1: a miss; then a 6 kills the fly")
	}
}
