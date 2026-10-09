package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLeaper(t *testing.T) {
	tb := slayTable(t, "leaper", engine.SituationPlayer{Character: "isaac", Items: items("breakfast")}, seat("cain"))
	tb.G.ForceRolls(1, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "leaper")
	if d := tb.G.Players[0].Damage; d != 2 {
		t.Errorf("damage = %d, want 2 (doubled on a roll of 1)", d)
	}
}
