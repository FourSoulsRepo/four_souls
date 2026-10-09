package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestFlush(t *testing.T) {
	tb := itemTable(t, items("flush"))
	fly, _ := tb.G.Monsters[0].TopOf()
	tb.Activate(0, "flush", 0, "Put each monster not being attacked on the bottom of the monster deck.", "leech")
	if tb.G.Object(fly).Zone.Kind == engine.ZoneInPlay {
		t.Error("the fly is still there")
	}
	deck := tb.G.Decks[engine.MonsterDeck]
	if a, b := tb.G.Object(deck[0]).Card, tb.G.Object(deck[1]).Card; a != "leech" || b != "fly" {
		t.Errorf("bottom of the monster deck is %s, %s; want the chosen leech lowest, then the fly", a, b)
	}
}
