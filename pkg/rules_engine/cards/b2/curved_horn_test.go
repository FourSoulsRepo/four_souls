package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestCurvedHorn(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "curved_horn"), seat("cain")},
		Monsters: []engine.CardRef{"conjoined_fatty"}, // 4 HP
	}, Set)
	tb.G.Turn.Attacks = 1
	fatty, _ := tb.G.Monsters[0].TopOf()
	tb.G.ForceRolls(6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "conjoined_fatty")
	// 2 damage on the first roll, then 1 and 1: dead after the third roll.
	if tb.G.Object(fatty).Zone.Kind == engine.ZoneInPlay {
		t.Error("the fatty survived 2 + 1 + 1 damage")
	}
	if tb.G.Turn.AttackRolls != 3 {
		t.Errorf("%d attack rolls, want 3", tb.G.Turn.AttackRolls)
	}
}
