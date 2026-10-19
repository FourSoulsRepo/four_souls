package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCharmedKeeper(t *testing.T) {
	tb := fight(t, "charmed_keeper")
	hurt(tb, 0, 1)
	tb.G.ForceRolls(6, 2, 4) // the kill; the reward roll 2; Cain's roll 4
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "charmed_keeper", foe)
	if tb.G.Players[0].Cents != 2 || tb.G.Players[1].Cents != 4 {
		t.Errorf("Isaac 2¢ from the reward roll, Cain 4¢; got %d and %d", tb.G.Players[0].Cents, tb.G.Players[1].Cents)
	}
}
