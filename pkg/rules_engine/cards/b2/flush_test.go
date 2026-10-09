package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestFlush(t *testing.T) {
	tb := itemTable(t, items("flush"))
	fly, _ := tb.G.Monsters[0].TopOf()
	tb.Activate(0, "flush", 0, "Put each monster not being attacked on the bottom of the monster deck.")
	if tb.G.Object(fly).Zone.Kind == engine.ZoneInPlay {
		t.Error("the fly is still there")
	}
	if c := tb.G.Object(tb.G.Decks[engine.MonsterDeck][0]).Card; c != "fly" && c != "leech" {
		t.Errorf("bottom of the monster deck is %s", c)
	}
}
