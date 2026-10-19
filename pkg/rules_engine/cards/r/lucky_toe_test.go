package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLuckyToe(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"lucky_toe"}, Hand: []engine.CardRef{"a_sack"}}, seat("cain"),
	}})
	tb.Play(0, "a_sack") // loot 3
	if got := len(tb.G.Players[0].Hand); got != 4 {
		t.Errorf("hand = %d, want 3 + 1", got)
	}
}

func TestLuckyToeNotInLootStep(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("cain"), {Character: "isaac", Items: []engine.CardRef{"lucky_toe"}},
	}})
	tb.EndTurn() // Isaac's turn: the loot step loots 1
	if got := len(tb.G.Players[1].Hand); got != 1 {
		t.Errorf("hand = %d, want 1: the loot step is not changed", got)
	}
}
