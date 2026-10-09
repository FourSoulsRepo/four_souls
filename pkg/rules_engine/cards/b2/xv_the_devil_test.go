package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestXVTheDevil(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			{Character: "isaac", Items: []engine.CardRef{"curved_horn"}, Hand: []engine.CardRef{"xv_the_devil"}},
			{Character: "cain", Items: []engine.CardRef{"swallowed_penny", "the_d6"}},
		},
	}, Set)
	if opts := len(tb.G.Players[1].InPlay); opts != 2 {
		t.Fatal("setup")
	}
	tb.Play(0, "xv_the_devil", "curved_horn", "swallowed_penny")
	if !hasItem(tb.G, 0, "swallowed_penny") || hasItem(tb.G, 0, "curved_horn") {
		t.Error("the horn was not traded for the penny")
	}
}
