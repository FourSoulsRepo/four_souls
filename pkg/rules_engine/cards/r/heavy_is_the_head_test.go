package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHeavyIsTheHead(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "jawbone"), seat("cain", "spoon_bender", "glitch")},
		Rooms:   []engine.CardRef{"heavy_is_the_head"},
	})
	for _, id := range tb.G.Players[1].InPlay {
		tb.G.Object(id).Charged = false
	}
	tb.EndTurn("spoon_bender")
	if tb.G.Object(tb.Find(1, "glitch")).Charged {
		t.Error("Cain controls the most items: he recharges only one")
	}
}
