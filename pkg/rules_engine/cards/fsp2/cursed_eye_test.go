package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestCursedEye(t *testing.T) {
	tb := armed(t, items("cursed_eye"), "monstro", "fatty")
	tb.G.ForceRolls(6) // 1 + 3 = 4 damage kills Monstro (HP 4)
	tb.Do(attackIntent(tb), "monstro")
	if top(tb, 0) == "monstro" {
		t.Error("a 6 deals 3 more")
	}
	tb2 := armed(t, items("cursed_eye"), "monstro", "fatty")
	tb2.G.ForceRolls(1)
	ev := tb2.Do(attackIntent(tb2), "monstro")
	if !enginetest.Has(ev, engine.EvTurnEndedEarly, 0) {
		t.Error("a 1 ends the turn")
	}
	if tb2.G.Players[0].Damage != 0 {
		t.Error("everything unresolved is cancelled: no combat damage from the 1")
	}
}
