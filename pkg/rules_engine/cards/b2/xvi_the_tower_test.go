package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestXVITheTower(t *testing.T) {
	tb := lootTable(t, "xvi_the_tower")
	fly, _ := tb.G.Monsters[0].TopOf()
	tb.G.ForceRolls(3)
	tb.Play(0, "xvi_the_tower", "the rest in slot order")
	if tb.G.Object(fly).Zone.Kind == engine.ZoneInPlay {
		t.Error("the fly survived 1 damage")
	}
	tb = lootTable(t, "xvi_the_tower")
	tb.G.ForceRolls(1)
	tb.Play(0, "xvi_the_tower")
	if a, b := tb.G.Players[0].Damage, tb.G.Players[1].Damage; a != 1 || b != 1 {
		t.Errorf("damage %d and %d, want 1 each", a, b)
	}
}
