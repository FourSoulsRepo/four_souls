package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHaloOfFlies(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "halo_of_flies"), seat("cain")}, Monsters: []engine.CardRef{"leech"}})
	tb.Attack("leech", 1, 6)
	if tb.G.Players[0].Damage != 0 {
		t.Error("no combat damage on an attack roll of 1")
	}
}

func TestHaloOfFliesWhenAttacked(t *testing.T) {
	tb := sirenTable(t, nil, []engine.CardRef{"halo_of_flies"})
	tb.G.ForceRolls(6, 4, 4)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe, "halo_of_flies")
	if tb.G.Players[1].Damage != 0 && !tb.G.Players[1].Dead {
		t.Error("a 6 deals no damage to Cain (R-ATK-24)")
	}
	if !tb.G.Players[1].Dead {
		t.Error("the 4s still hit")
	}
}
