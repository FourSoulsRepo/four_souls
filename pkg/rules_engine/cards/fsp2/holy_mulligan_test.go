package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHolyMulligan(t *testing.T) {
	tb := fight(t, "holy_mulligan", "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "dip", "spider", "trite", "leech") // the new slots fill with these
	slots := len(tb.G.Monsters)
	tb.G.ForceRolls(6)
	tb.Do(attackIntent(tb), "holy_mulligan")
	if n := len(tb.G.Monsters); n != slots+2 {
		t.Errorf("monster slots = %d, want %d", n, slots+2)
	}
	if tb.G.Turn.Attacks != 1 {
		t.Errorf("attacks left = %d, want 1 additional", tb.G.Turn.Attacks)
	}
}
