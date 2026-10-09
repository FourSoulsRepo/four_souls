package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheD10(t *testing.T) {
	tb := itemTable(t, items("the_d10", "mystery_sack"))
	fly, _ := tb.G.Monsters[0].TopOf()
	tb.G.ForceRolls(3)
	tb.Activate(0, "mystery_sack", 0, "slot 1")
	if top, _ := tb.G.Monsters[0].TopOf(); top == fly || tb.G.Object(fly).Zone.Kind != engine.ZoneCovered {
		t.Error("the fly was not covered")
	}
}
