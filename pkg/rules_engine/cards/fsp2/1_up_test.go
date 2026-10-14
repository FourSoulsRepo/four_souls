package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func Test1Up(t *testing.T) {
	tb := armed(t, items("1_up"), "leech", "fatty")
	tb.G.ForceRolls(1) // the leech hits Isaac for 2: lethal
	ev := tb.Do(attackIntent(tb), "leech")
	if tb.G.Players[0].Dead {
		t.Fatal("1-Up should prevent the death")
	}
	if tb.G.Players[0].Damage != 0 {
		t.Errorf("damage = %d: heal to full", tb.G.Players[0].Damage)
	}
	if hasItem(tb.G, 0, "1_up") {
		t.Error("1-Up is destroyed")
	}
	if tb.G.Attack.On {
		t.Error("the attack is cancelled")
	}
	if tb.G.Turn.Attacks != 1 {
		t.Errorf("attacks left = %d, want 1 additional", tb.G.Turn.Attacks)
	}
	// One popup for the whole heal (19.2): Isaac had 2 HP.
	var heals []int
	for _, e := range ev {
		if e.Kind == engine.EvHealed && e.Player == 0 {
			heals = append(heals, e.Amount)
		}
	}
	if len(heals) != 1 || heals[0] != 2 {
		t.Errorf("heals = %v, want one of 2", heals)
	}
}
