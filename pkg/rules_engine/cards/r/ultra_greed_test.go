package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestUltraGreedTeamUp(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"ultra_greed"},
	})
	tb.G.Object(tb.G.Monsters[0].Cards[0]).Damage = 7 // 1 HP left
	// Isaac misses (1 damage to him), Cain misses, Samson hits and kills.
	tb.G.ForceRolls(1, 1, 5)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "ultra_greed")
	if tb.G.Players[2].Cents != 40 || tb.G.SoulValue(2) != 1 {
		t.Errorf("Samson killed it: he gains the 40¢ and the soul; cents %d, souls %d", tb.G.Players[2].Cents, tb.G.SoulValue(2))
	}
	if tb.G.Players[0].Damage != 1 || tb.G.Players[1].Damage != 1 {
		t.Error("each miss hurts the player who rolled it")
	}
}
