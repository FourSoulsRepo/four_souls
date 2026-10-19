package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAbel(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "abel"), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"pooter", "fly"},
	})
	tb.G.ForceRolls(6, 5, 6, 1) // hit; Abel's 5; hit again; Abel's 1
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "pooter", "fly")
	if top(tb, 1) == "fly" {
		t.Error("Abel's 4-6 deals 1 damage to another monster: the fly dies")
	}
}
