package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCarrionQueen(t *testing.T) {
	tb := slayTable(t, "carrion_queen", seat("isaac"), seat("cain"))
	m, _ := tb.G.Monsters[0].TopOf()
	tb.G.ForceRolls(4, 5, 6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "carrion_queen")
	if tb.G.Object(m).Zone.Kind == engine.ZoneInPlay {
		t.Fatal("the queen survived three hits of 6")
	}
	if tb.G.Turn.AttackRolls != 5 {
		t.Errorf("%d attack rolls, want 5: rolls of 4 and 5 did nothing", tb.G.Turn.AttackRolls)
	}
}
