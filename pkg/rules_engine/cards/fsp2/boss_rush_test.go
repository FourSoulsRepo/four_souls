package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBossRush(t *testing.T) {
	tb := fight(t, "fatty", "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "widow", "dip", "monstro", "spider") // spider on top
	tb.RevealFromDeck("boss_rush", "slot 2", "slot 2")
	if top(tb, 1) != "widow" {
		t.Errorf("slot 2 = %s, want the second boss on top", top(tb, 1))
	}
	if tb.G.Turn.MustAttacks != 1 {
		t.Errorf("must attacks = %d, want 1", tb.G.Turn.MustAttacks)
	}
}
