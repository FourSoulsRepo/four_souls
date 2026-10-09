package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestEvilTwin(t *testing.T) {
	tb := slayTable(t, "evil_twin", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "evil_twin")
	if !tb.G.Players[1].Dead {
		t.Error("the left player did not take the damage")
	}
}
