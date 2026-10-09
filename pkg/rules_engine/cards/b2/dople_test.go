package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDople(t *testing.T) {
	tb := slayTable(t, "dople", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "dople")
	if d := tb.G.Players[1].Damage; d != 2 {
		t.Errorf("the right player's damage = %d, want 2", d)
	}
}
