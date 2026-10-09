package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheDukeOfFlies(t *testing.T) {
	tb := slayTable(t, "the_duke_of_flies", seat("isaac"), seat("cain"))
	m, _ := tb.G.Monsters[0].TopOf()
	// Each hit lets the active player roll: a 1 prevents it.
	tb.G.ForceRolls(6, 1, 6, 2, 6, 2, 6, 2, 6, 2)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "the_duke_of_flies")
	if tb.G.Object(m).Zone.Kind == engine.ZoneInPlay {
		t.Fatal("the duke survived")
	}
	if n := tb.G.Turn.AttackRolls; n != 5 {
		t.Errorf("%d attack rolls, want 5: the first hit was prevented", n)
	}
}
