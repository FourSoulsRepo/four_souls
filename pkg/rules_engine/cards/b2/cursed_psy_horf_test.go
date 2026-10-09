package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCursedPsyHorf(t *testing.T) {
	tb := slayTable(t, "cursed_psy_horf", engine.SituationPlayer{Character: "isaac", Items: items("mystery_sack")}, seat("cain"))
	tb.G.ForceRolls(6)
	tb.Activate(0, "mystery_sack", 0)
	if d := tb.G.Players[0].Damage; d != 1 {
		t.Errorf("damage = %d, want 1", d)
	}
}
