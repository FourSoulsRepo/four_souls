package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTechX(t *testing.T) {
	tb := itemTable(t, items("tech_x"))
	tb.G.Object(tb.Find(0, "tech_x")).Counters = []engine.Counter{{Count: 3}}
	tb.Activate(0, "tech_x", 1, foe)
	if !tb.G.Players[1].Dead {
		t.Error("not killed")
	}
}
