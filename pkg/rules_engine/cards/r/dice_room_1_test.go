package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDiceRoom1(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "spoon_bender"), seat("cain", "jawbone")},
	})
	enterRoom(tb, "dice_room_1")
	if hasItem(tb.G, 0, "spoon_bender") || !hasItem(tb.G, 1, "jawbone") {
		t.Error("only the active player rerolls his items")
	}
	tb.EndTurn()
	if len(tb.G.Rooms[0].Cards) != 1 || tb.G.Object(room(tb)).Card == "dice_room_1" {
		t.Error("it goes to discard at the end of the turn; the slot is refilled")
	}
}
