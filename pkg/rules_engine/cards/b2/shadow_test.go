package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestShadowTakesThePenalty(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: items("breakfast"), Hand: items("xiii_death", "a_dime"), Cents: 3}, seat("cain", "shadow"),
	}}, Set)
	tb.Play(0, "xiii_death", me, "breakfast", "a_dime") // Cain picks the item, Isaac the card
	if h := handCards(tb.G, 1); len(h) != 1 || h[0] != "a_dime" {
		t.Errorf("Cain's hand %v, want A Dime!!", h)
	}
	if c := tb.G.Players[1].Cents; c != 1 {
		t.Errorf("Cain's cents %d, want 1", c)
	}
	if hasItem(tb.G, 0, "breakfast") {
		t.Error("the breakfast was not destroyed")
	}
}
