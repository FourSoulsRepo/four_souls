package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLust(t *testing.T) {
	tb := slayTable(t, "lust", engine.SituationPlayer{Character: "isaac", Items: items("breakfast")}, seat("cain"))
	tb.G.ForceRolls(6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "lust")
	if d := tb.G.Players[0].Damage; d != 2 {
		t.Errorf("damage = %d, want 2", d)
	}
}
