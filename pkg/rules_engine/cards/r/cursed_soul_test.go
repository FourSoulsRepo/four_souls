package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCursedSoul(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("isaac", "cursed_soul")}})
	soul := tb.Find(1, "cursed_soul")
	tb.G.AddCountersTo(soul, 5)
	tb.EndTurn() // Isaac's turn: 1 damage, the 6th counter
	if tb.G.Players[1].Damage != 1 {
		t.Error("he takes 1 damage")
	}
	if tb.G.Object(soul).Role != engine.RoleSoul || tb.G.SoulValue(1) != 1 {
		t.Fatal("with 6 counters it becomes his soul")
	}
	tb.EndTurn()
	if tb.G.SoulValue(1) != 1 {
		t.Error("as a soul it has no abilities: it is not given away")
	}
}

func TestCursedSoulIsGivenAway(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "cursed_soul"), seat("cain")}})
	tb.EndTurn(foe)
	if !hasItem(tb.G, 1, "cursed_soul") {
		t.Error("at the end of his turn Isaac gives it to Cain")
	}
}
