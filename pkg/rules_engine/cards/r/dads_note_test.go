package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDadsNote(t *testing.T) {
	tb := lootTable(t, "dads_note")
	// Put two monsters and an event in the monster discard.
	for _, c := range []engine.CardRef{"pooter", "chest", "big_spider"} {
		tb.G.PutOnDeck(engine.MonsterDeck, c)
		tb.G.MillTop(engine.MonsterDeck)
	}
	tb.Play(0, "dads_note", "pooter", "ask.slot_card 1 fly", "big_spider", "ask.slot_card 1 fly", "ask.done")
	if got := top(tb, 0); got != "big_spider" {
		t.Errorf("slot 1 = %s, want big_spider on top", got)
	}
	if inDiscard(tb.G, engine.MonsterDeck, "pooter") || !inDiscard(tb.G, engine.MonsterDeck, "chest") {
		t.Error("only non-event monster cards come back")
	}
}
