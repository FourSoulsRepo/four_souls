package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestFriendlySack(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "friendly_sack"), seat("cain"), seat("samson")}})
	tb.G.PutOnDeck(engine.LootDeck, "bomb", "a_nickel", "dice_shard")
	tb.Activate(0, "friendly_sack", 0, "a_nickel", "player 3 (samson)")
	if !handHas(tb.G, 0, "a_nickel") || len(tb.G.Players[0].Hand) != 1 {
		t.Error("Isaac keeps the nickel")
	}
	if len(tb.G.Players[2].Hand) != 2 || len(tb.G.Players[1].Hand) != 0 {
		t.Error("Samson gets the other two")
	}
}
