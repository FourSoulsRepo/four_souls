package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPeepingFatty(t *testing.T) {
	tb := fight(t, "peeping_fatty")
	hurt(tb, 0, 1)
	tb.G.PutOnDeck(engine.LootDeck, "bomb", "a_nickel", "dice_shard", "butter_bean", "a_sack")
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "peeping_fatty", "bomb", "a_sack", "a_nickel", "dice_shard", "butter_bean")
	// The reward (loot 2) comes after the trigger: the first two are drawn.
	if !handHas(tb.G, 0, "bomb") || !handHas(tb.G, 0, "a_sack") {
		t.Error("the bomb and the sack were put on top, then looted")
	}
}
