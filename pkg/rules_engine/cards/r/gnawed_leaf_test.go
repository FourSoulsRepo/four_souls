package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGnawedLeafNotYourTurn(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"mega_bomb"}}, seat("isaac", "gnawed_leaf"), seat("samson"),
	}})
	tb.Play(0, "mega_bomb", "player 2 (isaac)", "No")
	if tb.G.Players[1].Damage != 0 {
		t.Error("damage is prevented while it's not his turn")
	}
}

func TestGnawedLeafCancels(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"two_of_spades", "two_of_spades"}},
		{Character: "isaac", Items: []engine.CardRef{"gnawed_leaf"}, Hand: []engine.CardRef{"bomb"}},
		seat("samson"),
	}})
	tb.G.Turn.LootPlays = 2
	tb.Play(0, "two_of_spades", "player 2 (isaac)", "yes")
	if len(tb.G.Players[0].Hand) != 1 {
		t.Fatal("the first loot aimed at him is cancelled")
	}
	tb.Play(0, "two_of_spades", "player 2 (isaac)")
	if len(tb.G.Players[0].Hand) != 1 {
		t.Error("the second resolves: Cain loots 1 for Isaac's 1 card")
	}
}
