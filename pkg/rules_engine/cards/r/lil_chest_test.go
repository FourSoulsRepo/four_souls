package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLilChest(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "lil_chest"), seat("cain")}, Monsters: []engine.CardRef{"leech"}})
	tb.Attack("leech", 6, 2) // the kill; Lil Chest's 2
	if got := tb.G.Players[0].Cents; got != 2 {
		t.Errorf("cents = %d, want 2", got)
	}
}
