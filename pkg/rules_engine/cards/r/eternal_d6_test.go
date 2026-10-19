package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestEternalD6(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "eternal_d6"), seat("cain")},
		Monsters: []engine.CardRef{"fly"},
	})
	tb.Activate(0, "eternal_d6", 0)
	tb.G.Object(tb.Find(0, "eternal_d6")).Charged = true
	tb.Activate(0, "eternal_d6", 0)
	tb.Activate(0, "eternal_d6", 1)
	tb.Attack("fly", 3, 6) // 3 becomes 1: a miss; then a hit
	if got := tb.G.Players[0].Damage; got != 1 {
		t.Errorf("damage = %d, want 1 from the miss", got)
	}
}
