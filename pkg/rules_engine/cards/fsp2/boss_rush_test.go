package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestBossRushOwnSlotInvariants: a boss put in Boss Rush's own slot
// leaves a consistent state (19.61). The slots and the discard are in
// testdata/situations/boss_rush/.
func TestBossRushOwnSlotInvariants(t *testing.T) {
	tb := fight(t, "fatty", "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "widow", "dip", "monstro", "spider") // spider on top
	tb.RevealFromDeck("boss_rush", "ask.slot 1", "ask.slot 2")
	if err := tb.G.CheckInvariants(); err != nil {
		t.Error(err)
	}
}
