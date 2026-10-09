package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestChest(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac"), seat("cain"))
	tb.G.ForceRolls(5)
	tb.RevealFromDeck("chest")
	if c := tb.G.Players[0].Cents; c != 6 {
		t.Errorf("cents = %d, want 6", c)
	}
	if !slices.ContainsFunc(tb.G.Discards[engine.MonsterDeck], func(id engine.ObjectID) bool { return tb.G.Object(id).Card == "chest" }) {
		t.Error("the event is not in the monster discard")
	}
	if top, _ := tb.G.Monsters[0].TopOf(); tb.G.Object(top).Card != "fly" {
		t.Error("the fly is not uncovered")
	}
}
