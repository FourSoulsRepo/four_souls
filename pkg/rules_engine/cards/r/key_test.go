package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestKeyAttackAgain(t *testing.T) {
	tb := lootTable(t, "key")
	tb.Play(0, "key", "The active player may attack an additional time this turn.")
	if got := tb.G.Turn.Attacks; got != 2 {
		t.Errorf("attacks = %d, want 2", got)
	}
}

func TestKeyDeckTopToSlot(t *testing.T) {
	tb := lootTable(t, "key")
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter")
	tb.Play(0, "key", "Put the top card of the monster deck in a monster slot not being attacked.", "slot 2: leech")
	if got := top(tb, 1); got != "pooter" {
		t.Errorf("slot 2 = %s, want pooter covering the leech", got)
	}
}
