package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCursedFatty(t *testing.T) {
	tb := slayTable(t, "cursed_fatty", engine.SituationPlayer{Character: "isaac", Items: items("mystery_sack"), Hand: items("a_dime")}, seat("cain"))
	tb.G.ForceRolls(5)
	tb.Activate(0, "mystery_sack", 0, "a_dime")
	if !inDiscard(tb.G, "a_dime") {
		t.Error("no discard on a 5")
	}
}
