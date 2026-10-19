package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCharmedMonstro(t *testing.T) {
	tb := fight(t, "charmed_monstro")
	hurt(tb, 0, 1)
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "charmed_monstro", "player 3 (samson)")
	if tb.G.SoulValue(0) != 0 || tb.G.SoulValue(2) != 1 || len(tb.G.Players[0].InPlay) != 2 {
		t.Error("Isaac gets the 2 treasures; Samson gets the soul")
	}
}
