package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSwarmOfFlies(t *testing.T) {
	tb := slayTable(t, "swarm_of_flies", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(5, 6, 6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "swarm_of_flies")
	if d := tb.G.Players[0].Damage; d != 1 {
		t.Errorf("damage = %d, want 1", d)
	}
}
