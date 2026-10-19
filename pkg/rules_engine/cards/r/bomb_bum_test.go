package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBombBum(t *testing.T) {
	tb := lootRoomTable(t, "bomb_bum", "bomb", "a_penny")
	tb.G.ForceRolls(1, 2, 5) // the last roll resolves first
	useRoom(tb, 0, 0, "bomb", foe, "leech")
	if tb.G.Players[0].Damage != 1 || tb.G.Players[1].Damage != 1 || !inDiscard(tb.G, engine.MonsterDeck, "leech") {
		t.Error("a bomb rolls 3 times: 5 hurts Isaac, 2 hurts Cain, 1 kills the leech")
	}
}
