package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAbaddon(t *testing.T) {
	tb := armed(t, items("abaddon"), "fatty", "monstro")
	if atk := tb.G.PlayerATK(0); atk != 4 {
		t.Errorf("ATK = %d, want 1 + 3", atk)
	}
	tb.G.DealDamageTo(engine.Target{Player: 0, IsPlayer: true}, 1, 1, 0)
	settleAny(tb)
	if !tb.G.Players[0].Dead {
		t.Error("any damage kills the Abaddon player")
	}
}
