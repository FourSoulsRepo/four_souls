package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCompost(t *testing.T) {
	tb := itemTable(t, items("compost", "bum_friend"), "a_dime")
	tb.Activate(0, "compost", 0)
	tb.G.ForceRolls()
	// Put a card into the loot discard, then loot through Bum Friend.
	tb.Play(0, "a_dime")
	tb.Activate(0, "bum_friend", 0, "a_dime")
	if top := tb.G.DeckTop(engine.LootDeck, 1); tb.G.Object(top[0]).Card != "a_dime" {
		t.Error("the looted card was not A Dime!! from the discard")
	}
}
