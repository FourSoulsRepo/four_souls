package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheScourgeWhip(t *testing.T) {
	tb := fight(t, "the_scourge")
	tb.G.ForceRolls(6, 1, 1)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "the_scourge", "whip", "whip", "whiff")
	if tb.G.Players[1].Damage != 1 || tb.G.Players[2].Damage != 1 {
		t.Error("whip wins: the damage is prevented, each non-active player takes 1")
	}
}

func TestTheScourgeWhiff(t *testing.T) {
	tb := fight(t, "the_scourge")
	hurt(tb, 0, 1)
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "the_scourge", "whiff", "whiff", "whip")
	if tb.G.SoulValue(0) != 1 || len(tb.G.Players[0].Hand) != 1+2 {
		t.Error("whiff wins: the damage kills it; loot 1, then its 2 rewards")
	}
}
