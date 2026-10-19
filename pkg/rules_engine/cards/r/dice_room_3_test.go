package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDiceRoom3(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac"), seat("cain")},
		Shop:    []engine.CardRef{"jawbone", "spoon_bender"},
	})
	enterRoom(tb, "dice_room_3")
	if !inDiscard(tb.G, engine.TreasureDeck, "jawbone") || !inDiscard(tb.G, engine.TreasureDeck, "spoon_bender") {
		t.Error("each shop item goes to discard")
	}
}
