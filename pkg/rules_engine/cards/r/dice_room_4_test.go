package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDiceRoom4(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "spoon_bender"), seat("cain", "jawbone")},
	})
	enterRoom(tb, "dice_room_4")
	if hasItem(tb.G, 0, "spoon_bender") || hasItem(tb.G, 1, "jawbone") {
		t.Error("each player rerolls each of their items")
	}
}
