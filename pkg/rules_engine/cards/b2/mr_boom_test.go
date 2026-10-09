package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMrBoom(t *testing.T) {
	tb := itemTable(t, items("mr_boom"))
	fly, _ := tb.G.Monsters[0].TopOf()
	tb.Activate(0, "mr_boom", 0, "fly")
	if tb.G.Object(fly).Zone.Kind == engine.ZoneInPlay {
		t.Error("the fly survived")
	}
}
