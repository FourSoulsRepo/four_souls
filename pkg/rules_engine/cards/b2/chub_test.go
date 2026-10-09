package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestChub(t *testing.T) {
	tb := slayTable(t, "chub", seat("isaac"), seat("cain"))
	m, _ := tb.G.Monsters[0].TopOf()
	// Two hits (2 damage), a 1 (it heals 2, Isaac takes 1), then 4 hits.
	tb.G.ForceRolls(6, 6, 1, 6, 6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "chub")
	if tb.G.Object(m).Zone.Kind == engine.ZoneInPlay {
		t.Fatal("Chub survived")
	}
	if n := tb.G.Turn.AttackRolls; n != 7 {
		t.Errorf("%d attack rolls, want 7 (it healed 2)", n)
	}
}
