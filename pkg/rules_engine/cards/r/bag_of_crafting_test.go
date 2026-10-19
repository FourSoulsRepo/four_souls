package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBagOfCrafting(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "the_hoarder", Items: []engine.CardRef{"bag_of_crafting"}, Hand: []engine.CardRef{"bomb", "bomb", "bomb", "bomb"}}, seat("cain"),
	}})
	for range 4 {
		tb.Activate(0, "bag_of_crafting", 0, "bomb")
	}
	tb.Activate(0, "bag_of_crafting", 1)
	if len(tb.G.Players[0].InPlay) != 2 || len(tb.G.Players[0].Hand) != 0 {
		t.Error("4 discards, 4 counters: +1 treasure")
	}
}
