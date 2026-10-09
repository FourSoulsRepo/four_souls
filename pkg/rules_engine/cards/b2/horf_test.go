package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHorf(t *testing.T) {
	tb := slayTable(t, "horf", engine.SituationPlayer{Character: "isaac", Items: items("breakfast")}, seat("cain"))
	tb.G.ForceRolls(2, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "horf")
	if d := tb.G.Players[0].Damage; d != 2 {
		t.Errorf("damage = %d, want 2 (1 + 1 on a roll of 2)", d)
	}
}
