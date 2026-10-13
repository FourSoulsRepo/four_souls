package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBagOTrash(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "bum_bo_fsp2", Items: []engine.CardRef{"bag_o_trash"}, Cents: 13}, seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.Activate(0, "bag_o_trash", 0)
	if tb.G.Turn.LootPlays != 2 {
		t.Errorf("loot plays = %d, want 2", tb.G.Turn.LootPlays)
	}
	tb.Activate(0, "bag_o_trash", 1)
	if n := len(tb.G.Players[0].Hand); n != 1 {
		t.Errorf("hand = %d, want 1", n)
	}
	tb.Activate(0, "bag_o_trash", 2, "fly")
	if top(tb, 0) == "fly" {
		t.Error("the fly should die")
	}
	if c := tb.G.Players[0].Cents; c != 1 {
		t.Errorf("cents = %d, want 13 - 3 - 4 - 6 + the fly's 1¢", c)
	}
}
