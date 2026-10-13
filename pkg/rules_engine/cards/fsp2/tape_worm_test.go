package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTapeWorm(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"tape_worm"}}, seat("cain")},
		Monsters: []engine.CardRef{"fatty", "monstro"},
	})
	tb.G.ForceRolls(1, 2, 2, 2, 2)
	tb.Do(attackIntent(tb), "fatty", foe)
	if d := tb.G.Players[1].Damage; d != 1 {
		t.Errorf("Cain's damage = %d, want 1 from the miss", d)
	}
}
