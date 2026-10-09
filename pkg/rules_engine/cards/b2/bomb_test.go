package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBomb(t *testing.T) {
	tb := lootTable(t, "bomb")
	fly, _ := tb.G.Monsters[0].TopOf()
	tb.Play(0, "bomb", "fly")
	if tb.G.Object(fly).Zone.Kind == engine.ZoneInPlay {
		t.Error("the fly survived 1 damage")
	}
}
