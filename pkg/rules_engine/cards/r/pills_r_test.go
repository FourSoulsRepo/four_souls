package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPillsLootThree(t *testing.T) {
	tb := lootTable(t, "pills_r")
	tb.G.ForceRolls(6)
	tb.Play(0, "pills_r")
	if got := len(tb.G.Players[0].Hand); got != 3 {
		t.Errorf("hand = %d, want 3", got)
	}
}

func TestPillsPutOneBack(t *testing.T) {
	tb := lootTable(t, "pills_r", "bomb")
	tb.G.ForceRolls(3)
	tb.Play(0, "pills_r", "bomb")
	if got := len(tb.G.Players[0].Hand); got != 3 {
		t.Errorf("hand = %d, want 1 + 3 - 1", got)
	}
	if top := tb.G.DeckTop(engine.LootDeck, 1); len(top) == 0 || tb.G.Object(top[0]).Card != "bomb" {
		t.Error("the bomb should be on top of the loot deck")
	}
}

func TestPillsPutThreeBack(t *testing.T) {
	tb := lootTable(t, "pills_r", "bomb", "a_nickel", "dice_shard")
	tb.G.ForceRolls(1)
	tb.Play(0, "pills_r", "bomb", "a_nickel", "dice_shard")
	if got := len(tb.G.Players[0].Hand); got != 3 {
		t.Errorf("hand = %d, want 3 + 3 - 3", got)
	}
	var got []engine.CardRef
	for _, id := range tb.G.DeckTop(engine.LootDeck, 3) {
		got = append(got, tb.G.Object(id).Card)
	}
	if len(got) != 3 || got[0] != "dice_shard" || got[2] != "bomb" {
		t.Errorf("deck top = %v, want the last one picked on top", got)
	}
}
