package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBrimstone(t *testing.T) {
	tb := monsterTable(t, items("conjoined_fatty"), seat("isaac", "brimstone"), seat("cain"))
	if a := tb.G.PlayerATK(0); a != 2 {
		t.Errorf("ATK = %d, want 2", a)
	}
	tb.G.ForceRolls(6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "conjoined_fatty", foe, foe)
	if d := tb.G.Players[1].Damage; d != 2 {
		t.Errorf("Cain's damage = %d, want 2 (one per combat damage)", d)
	}
}
