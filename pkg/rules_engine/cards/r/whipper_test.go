package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestWhipper(t *testing.T) {
	tb := fight(t, "whipper")
	tb.G.ForceRolls(2, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly")
	if tb.G.Players[0].Damage != 1 {
		t.Error("the declaration's roll of 2: Whipper deals 1 damage to Isaac")
	}
}
