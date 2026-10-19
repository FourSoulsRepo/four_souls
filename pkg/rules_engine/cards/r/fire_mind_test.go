package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestFireMind(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "fire_mind"), seat("cain")}, Monsters: []engine.CardRef{"pooter"}})
	tb.Attack("pooter", 6) // 1 ATK doubled kills the 2-HP pooter
	if top(tb, 0) == "pooter" {
		t.Error("damage doubled on a 6")
	}
	tb2 := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "fire_mind", "spoon_bender"), seat("cain")}, Monsters: []engine.CardRef{"pooter"}})
	tb2.G.ForceRolls(1)
	tb2.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "pooter", "spoon_bender")
	if !tb2.G.Players[0].Dead {
		t.Error("damage taken doubled on a 1: 2 damage kills Isaac")
	}
}
