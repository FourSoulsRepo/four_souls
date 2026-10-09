package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMomsHand(t *testing.T) {
	tb := slayTable(t, "moms_hand", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "moms_hand")
	if !tb.G.Turn.EndDeclared && tb.G.Turn.Active == 0 {
		t.Error("the turn did not end")
	}
}
