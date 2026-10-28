package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHolyGreedling(t *testing.T) {
	tb := fight(t, "holy_greedling")
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "ask.yes")
	if tb.G.Players[0].Cents != 1+3 {
		t.Errorf("the fly's 1¢ and Holy Greedling's 3¢: cents %d", tb.G.Players[0].Cents)
	}
}

// The refill waits for Holy Greedling's reward roll (19.104): Holy Dinga,
// the next monster, must not see the 6 of that roll.
func TestHolyGreedlingRefillWaitsForTheRewardRoll(t *testing.T) {
	tb := fight(t, "holy_greedling")
	hurt(tb, 0, 1)
	tb.G.PutOnDeck(engine.MonsterDeck, "holy_dinga")
	tb.G.Players[0].Damage = 1
	tb.G.ForceRolls(6, 6) // the attack kills it; the reward roll is a 6
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "holy_greedling", "ask.yes")
	if tb.G.Players[0].Damage != 1 {
		t.Errorf("Holy Dinga healed on the reward roll: damage %d, want 1", tb.G.Players[0].Damage)
	}
}
