package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestFlush(t *testing.T) {
	tb := itemTable(t, items("flush"))
	fly, _ := tb.G.Monsters[0].TopOf()
	leech, _ := tb.G.Monsters[1].TopOf()
	tb.Activate(0, "flush", 0, "Put each monster not being attacked on the bottom of the monster deck.")
	if tb.G.Object(fly).Zone.Kind == engine.ZoneInPlay || tb.G.Object(leech).Zone.Kind == engine.ZoneInPlay {
		t.Error("a monster is still there")
	}
	// No order question: random order at the bottom (Q-22).
	deck := tb.G.Decks[engine.MonsterDeck]
	got := []engine.CardRef{tb.G.Object(deck[0]).Card, tb.G.Object(deck[1]).Card}
	if !slices.Contains(got, "fly") || !slices.Contains(got, "leech") {
		t.Errorf("bottom of the monster deck is %v; want the fly and the leech", got)
	}
}
