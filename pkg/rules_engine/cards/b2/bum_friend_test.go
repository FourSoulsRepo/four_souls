package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBumFriend(t *testing.T) {
	tb := itemTable(t, items("bum_friend"), "a_dime")
	tb.Activate(0, "bum_friend", 0, "a_dime")
	if len(tb.G.Players[0].Hand) != 1 || tb.G.Object(tb.G.DeckTop(engine.LootDeck, 1)[0]).Card != "a_dime" {
		t.Error("A Dime!! is not back on top")
	}
}
