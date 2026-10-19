package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestUndefined(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "undefined"), seat("cain")},
		Shop:    []engine.CardRef{"sack_of_pennies", "spoon_bender"},
	})
	tb.Activate(0, "undefined", 0, "sack_of_pennies")
	id := tb.Find(0, "undefined")
	if tb.G.CardOf(id) != "sack_of_pennies" || !tb.G.Object(id).Charged {
		t.Fatal("it becomes a copy of the sack and recharges")
	}
	tb.EndTurn()
	if tb.G.CardOf(id) != "sack_of_pennies" {
		t.Error("the copy lasts past the end of the turn")
	}
	tb.EndTurn("spoon_bender") // Isaac's turn: the copy ends; he may discard a shop item
	if tb.G.CardOf(id) != "undefined" {
		t.Error("the copy ends at the start of his next turn")
	}
	if !inDiscard(tb.G, engine.TreasureDeck, "spoon_bender") {
		t.Error("he put the Spoon Bender into discard")
	}
}
