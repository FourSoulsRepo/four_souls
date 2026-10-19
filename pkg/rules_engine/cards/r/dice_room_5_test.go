package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDiceRoom5(t *testing.T) {
	tb := roomTable(t, "")
	enterRoom(tb, "dice_room_5")
	if !inDiscard(tb.G, engine.MonsterDeck, "fly") || !inDiscard(tb.G, engine.MonsterDeck, "leech") {
		t.Error("each monster goes to discard")
	}
}
