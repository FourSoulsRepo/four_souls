package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestEdensBlessingR(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac"), seat("cain", "spoon_bender")},
		Rooms:   []engine.CardRef{"edens_blessing_r"},
	})
	tb.EndTurn("spoon_bender")
	if hasItem(tb.G, 1, "spoon_bender") || !hasItem(tb.G, 1, "book_of_sin") {
		t.Error("Cain may reroll an item at the start of his turn")
	}
}
