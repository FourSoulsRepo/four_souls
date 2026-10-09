package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDaddyLongLegs(t *testing.T) {
	tb := slayTable(t, "daddy_long_legs", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(1, 6, 6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "daddy_long_legs")
	if len(tb.G.Boosts) == 0 && tb.G.Turn.Step == engine.StepAction {
		t.Error("no +1 DC boost on a roll of 1")
	}
}
