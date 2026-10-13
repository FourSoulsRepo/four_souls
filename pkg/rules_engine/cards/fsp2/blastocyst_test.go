package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBlastocyst(t *testing.T) {
	tb := fight(t, "blastocyst", "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "dip", "spider", "trite", "leech") // the new slots fill with these
	slots := len(tb.G.Monsters)
	tb.G.ForceRolls(6, 6, 6, 6, 6)
	tb.Do(attackIntent(tb), "blastocyst")
	if n := len(tb.G.Monsters); n != slots+2 {
		t.Errorf("monster slots = %d, want %d", n, slots+2)
	}
	if tb.G.Turn.DeckAttacks != 1 {
		t.Errorf("deck attacks = %d, want 1", tb.G.Turn.DeckAttacks)
	}
}
