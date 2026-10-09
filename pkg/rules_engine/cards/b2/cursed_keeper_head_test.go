package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCursedKeeperHead(t *testing.T) {
	tb := slayTable(t, "cursed_keeper_head", engine.SituationPlayer{Character: "isaac", Items: items("mystery_sack"), Cents: 3}, seat("cain"))
	tb.G.ForceRolls(1)
	tb.Activate(0, "mystery_sack", 0)
	if c := tb.G.Players[0].Cents; c != 1 {
		t.Errorf("cents = %d, want 1", c)
	}
}
