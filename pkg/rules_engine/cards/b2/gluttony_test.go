package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGluttony(t *testing.T) {
	tb := slayTable(t, "gluttony", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6, 6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "gluttony")
	if !tb.G.Players[1].Dead {
		t.Errorf("Cain's damage = %d, want 2 or dead", tb.G.Players[1].Damage)
	}
}
