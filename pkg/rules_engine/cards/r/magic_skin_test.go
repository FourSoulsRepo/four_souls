package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMagicSkin(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "magic_skin"), seat("cain", "spoon_bender")}})
	tb.Activate(0, "magic_skin", 0, foe)
	if len(tb.G.Players[0].InPlay) != 1 || !hasItem(tb.G, 1, "magic_skin") {
		t.Fatal("Isaac gains a treasure and gives the skin to Cain")
	}
	tb.G.ForceRolls(2)
	tb.EndTurn() // Cain's turn: a roll of 2 is not less than 1 counter
	if !hasItem(tb.G, 1, "magic_skin") || !hasItem(tb.G, 1, "spoon_bender") {
		t.Error("2 >= 1: nothing is destroyed")
	}
}

func TestMagicSkinDestroys(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain", "the_d6"), seat("isaac", "magic_skin", "spoon_bender")}})
	tb.G.AddCountersTo(tb.Find(1, "magic_skin"), 3)
	tb.G.ForceRolls(2)
	tb.EndTurn()
	if hasItem(tb.G, 1, "magic_skin") || hasItem(tb.G, 1, "spoon_bender") {
		t.Error("2 < 3: the skin and his other items are destroyed")
	}
}
