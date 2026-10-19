package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBetrayal(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), {Character: "cain", Cents: 7, Souls: []engine.CardRef{"monstro"}}},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.ForceRolls(6, 6)
	tb.G.PutOnDeck(engine.MonsterDeck, "betrayal")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", foe, "monstro")
	if tb.G.SoulValue(0) != 1 || tb.G.SoulValue(1) != 0 {
		t.Error("6: Isaac steals Cain's soul")
	}
}
