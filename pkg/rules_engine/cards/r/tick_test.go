package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTick(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"tick"}, Damage: 1}, seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.Attack("fly", 6)
	if got := tb.G.Players[0].Damage; got != 0 {
		t.Errorf("damage = %d, want healed by 1", got)
	}
	tb.Attack("leech", 6) // the additional attack
	if got := tb.G.Players[0].Cents; got != 1 {
		t.Errorf("cents = %d, want the fly's 1", got)
	}
	if tb.G.Turn.Attacks != 0 {
		t.Errorf("attacks left = %d; the second kill gives no more", tb.G.Turn.Attacks)
	}
}
