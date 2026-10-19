package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLaserEye(t *testing.T) {
	tb := lootRoomTable(t, "laser_eye")
	tb.G.PutOnDeck(engine.LootDeck, "a_penny")
	tb.G.PutOnDeck(engine.LootDeck, "a_penny")
	tb.EndTurn()
	if got := tb.G.Players[1].Damage; got != 1 {
		t.Errorf("Cain's damage after his loot step = %d, want 1", got)
	}
}
