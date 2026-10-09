package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestEhwaz(t *testing.T) {
	tb := lootTable(t, "ehwaz")
	fly, _ := tb.G.Monsters[0].TopOf()
	leech, _ := tb.G.Monsters[1].TopOf()
	tb.Play(0, "ehwaz")
	for _, id := range []engine.ObjectID{fly, leech} {
		if tb.G.Object(id).Zone.Kind == engine.ZoneInPlay {
			t.Errorf("%s is still in play", tb.G.Object(id).Card)
		}
	}
	for i, s := range tb.G.Monsters {
		if _, ok := s.TopOf(); !ok {
			t.Errorf("slot %d is empty", i)
		}
	}
}
