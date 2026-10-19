package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDiceRoom6(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"spoon_bender"}, Hand: []engine.CardRef{"a_penny"}}, seat("cain")},
		Monsters: []engine.CardRef{"fly", "leech"},
		Shop:     []engine.CardRef{"jawbone", "glitch"},
	})
	enterRoom(tb, "dice_room_6")
	if hasItem(tb.G, 0, "spoon_bender") || !inDiscard(tb.G, engine.TreasureDeck, "jawbone") {
		t.Error("each item in play is rerolled")
	}
	if handHas(tb.G, 0, "a_penny") || len(tb.G.Players[0].Hand) != 3 {
		t.Error("each player discards their hand and loots 3")
	}
	if !inDiscard(tb.G, engine.MonsterDeck, "fly") {
		t.Error("each monster goes to discard")
	}
}
