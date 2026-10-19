package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHolyGreedling(t *testing.T) {
	tb := fight(t, "holy_greedling")
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "yes")
	if tb.G.Players[0].Cents != 1+3 {
		t.Errorf("the fly's 1¢ and Holy Greedling's 3¢: cents %d", tb.G.Players[0].Cents)
	}
}
