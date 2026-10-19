package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestRibOfGreedDeathPenalty(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"mega_bomb"}},
		{Character: "cain", Items: []engine.CardRef{"rib_of_greed", "spoon_bender"}, Hand: []engine.CardRef{"a_nickel"}, Cents: 3},
	}})
	tb.Play(0, "mega_bomb", foe, "spoon_bender") // the penalty item: the Spoon Bender
	if !tb.G.Players[1].Dead {
		t.Fatal("Cain should be dead")
	}
	if got := tb.G.Players[1].Cents; got != 3 {
		t.Errorf("cents = %d, want 3: the penalty can't take them", got)
	}
	if got := len(tb.G.Players[1].Hand); got != 1 {
		t.Errorf("hand = %d, want 1: the penalty can't discard", got)
	}
}

func TestRibOfGreedAllowsSteals(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"jawbone"}}, {Character: "cain", Items: []engine.CardRef{"rib_of_greed"}, Cents: 5},
	}})
	tb.Activate(0, "jawbone", 0, foe)
	if got := tb.G.Players[1].Cents; got != 2 {
		t.Errorf("Cain's cents = %d, want 5 - 3: a steal is a give (Q-17)", got)
	}
}

func TestRibOfGreedAbilityDiscard(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"rib_of_greed"}, Hand: []engine.CardRef{"pills_r", "bomb"}}, seat("cain"),
	}})
	tb.G.ForceRolls(3)
	tb.Play(0, "pills_r", "bomb") // putting a card on the deck is not a discard
	if got := len(tb.G.Players[0].Hand); got != 3 {
		t.Errorf("hand = %d, want 1 + 3 - 1", got)
	}
}
