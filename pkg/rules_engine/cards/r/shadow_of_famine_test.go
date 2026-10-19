package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestShadowOfFamine(t *testing.T) {
	tb := lootRoomTable(t, "shadow_of_famine", "a_penny", "a_penny")
	tb.G.Turn.LootPlays = 2
	tb.Play(0, "a_penny")
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.G.Players[0].Hand[0]}}); err == nil {
		t.Error("a second loot card this turn")
	}
	tb.Activate(0, "isaac", 0)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{room(tb)}}); err == nil {
		t.Error("a second ability this turn")
	}
}
