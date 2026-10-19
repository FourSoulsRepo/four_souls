package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHemoptysis(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("the_benighted", "hemoptysis"), seat("cain")},
		Monsters: []engine.CardRef{"leech", "fly"},
	})
	h := tb.Find(0, "hemoptysis")
	tb.G.AddCountersTo(h, 3)
	tb.Activate(0, "hemoptysis", 1, "leech")
	if tb.G.MonsterATK(tb.G.Monsters[0].Cards[0]) != 1 {
		t.Error("the leech's ATK becomes 1")
	}
	tb.Activate(0, "hemoptysis", 2, "leech", "fly")
	if top(tb, 0) == "leech" || top(tb, 1) == "fly" {
		t.Error("1 damage to each of 2 monsters")
	}
}
