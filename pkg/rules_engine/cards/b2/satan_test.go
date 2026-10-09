package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSatan(t *testing.T) {
	tb := slayTable(t, "satan", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "satan", foe)
	if !tb.G.Players[1].Dead {
		t.Error("Cain did not die")
	}
}
