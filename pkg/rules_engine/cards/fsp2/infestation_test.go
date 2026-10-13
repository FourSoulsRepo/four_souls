package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestInfestation(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "guppy", Items: []engine.CardRef{"infestation"}, Hand: []engine.CardRef{"bomb"}}, seat("cain"),
	}})
	tb.Activate(0, "infestation", 0, "bomb")
	if n := len(tb.G.Players[0].Hand); n != 2 {
		t.Errorf("hand = %d, want 1 + 2 - 1", n)
	}
	if !inDiscard(tb.G, engine.LootDeck, "bomb") {
		t.Error("the bomb was not discarded")
	}
}
