package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestKeeperHead(t *testing.T) {
	tb := slayTable(t, "keeper_head", engine.SituationPlayer{Character: "isaac", Cents: 3}, seat("cain"))
	tb.G.ForceRolls(1, 6, 6, 1)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "keeper_head")
	if c := tb.G.Players[0].Cents; c != 1+1 {
		t.Errorf("cents = %d, want 2 (3 - 2, then a reward roll of 1)", c)
	}
}
