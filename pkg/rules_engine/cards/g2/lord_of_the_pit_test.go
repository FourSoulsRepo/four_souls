package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLordOfThePit(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("azazel", "lord_of_the_pit"), seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.ForceRolls(6)
	declareAttack(tb, "fly")
	// The attack roll is on the stack; cancel the attack in response.
	tb.Activate(0, "lord_of_the_pit", 0)
	if tb.G.Attack.On {
		t.Fatal("the attack was not cancelled")
	}
	fly, _ := tb.G.Monsters[0].TopOf()
	if tb.G.Object(fly).Card != "fly" {
		t.Error("the fly died: the cancelled attack roll must not resolve")
	}
	if tb.G.Turn.Attacks != 1 {
		t.Errorf("attacks left = %d, want 1 (an additional attack)", tb.G.Turn.Attacks)
	}
	// Without an attack it does nothing.
	tb2 := table(t, seat("azazel", "lord_of_the_pit"), seat("cain"))
	tb2.Activate(0, "lord_of_the_pit", 0)
	if tb2.G.Turn.Attacks != 1 {
		t.Errorf("attacks = %d, want 1: no attack was cancelled", tb2.G.Turn.Attacks)
	}
}
