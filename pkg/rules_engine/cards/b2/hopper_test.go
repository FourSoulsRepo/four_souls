package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHopper(t *testing.T) {
	tb := slayTable(t, "hopper", seat("isaac"), seat("cain"))
	m, _ := tb.G.Monsters[0].TopOf()
	tb.G.ForceRolls(6, 5, 5)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "hopper")
	if tb.G.Object(m).Zone.Kind == engine.ZoneInPlay {
		t.Fatal("it survived")
	}
	if tb.G.Turn.AttackRolls != 3 {
		t.Errorf("%d attack rolls, want 3: the 6 did no damage", tb.G.Turn.AttackRolls)
	}
}
