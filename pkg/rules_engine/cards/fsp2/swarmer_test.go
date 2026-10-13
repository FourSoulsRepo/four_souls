package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSwarmer(t *testing.T) {
	tb := fight(t, "swarmer", "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "dip", "spider", "trite", "leech") // the new slots fill with these
	slots := len(tb.G.Monsters)
	tb.G.KillObject(topID(tb, 0))
	tb.Settle()
	if n := len(tb.G.Monsters); n != slots+2 {
		t.Errorf("monster slots = %d, want %d", n, slots+2)
	}
}
