package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheBloat(t *testing.T) {
	tb := slayTable(t, "the_bloat", engine.SituationPlayer{Character: "isaac", Items: items("breakfast", "dinner")}, seat("cain"))
	tb.G.ForceRolls(1, 6, 6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "the_bloat")
	if d := tb.G.Players[1].Damage; d != 1 {
		t.Errorf("Cain's damage = %d, want 1", d)
	}
}
