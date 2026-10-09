package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSpiderMod(t *testing.T) {
	tb := itemTable(t, items("spider_mod", "mystery_sack"))
	fly, _ := tb.G.Monsters[0].TopOf()
	tb.G.ForceRolls(5)
	tb.Activate(0, "mystery_sack", 0, "fly")
	if tb.G.Object(fly).Zone.Kind == engine.ZoneInPlay {
		t.Error("the fly is still there")
	}
}
