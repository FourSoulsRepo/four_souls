package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestScolex(t *testing.T) {
	tb := slayTable(t, "scolex", engine.SituationPlayer{Character: "isaac", Items: items("breakfast"), Hand: items("a_dime")}, seat("cain"))
	tb.G.ForceRolls(1, 6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "scolex", "a_dime")
	if !inDiscard(tb.G, "a_dime") {
		t.Error("A Dime!! was not discarded")
	}
}
