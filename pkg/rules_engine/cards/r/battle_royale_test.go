package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBattleRoyale(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain", "spoon_bender")},
		Monsters: []engine.CardRef{"fly", "leech"},
		Rooms:    []engine.CardRef{"battle_royale"},
	})
	tb.G.ForceRolls(3, 3)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe, "spoon_bender")
	if !tb.G.Players[1].Dead || !hasItem(tb.G, 0, "spoon_bender") {
		t.Error("Cain (3+) dies and gives Isaac the item he would destroy")
	}
}
