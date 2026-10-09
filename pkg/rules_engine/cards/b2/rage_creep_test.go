package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestRageCreep(t *testing.T) {
	tb := slayTable(t, "rage_creep", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(1, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "rage_creep")
	if d := tb.G.Players[1].Damage; d != 1 {
		t.Errorf("Cain's damage = %d, want 1", d)
	}
}
