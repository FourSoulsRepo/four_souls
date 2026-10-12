package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestInfamy(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"infamy"}}, seat("cain")},
		Monsters: []engine.CardRef{"leech", "fly"},
	})
	tb.Attack("leech", 1, 6)
	if d := tb.G.Players[0].Damage; d != 0 {
		t.Errorf("damage = %d, want 0 on an attack roll of 1", d)
	}
	// Other misses still hurt: a 3 against the trite (DC 5) deals 1.
	tb2 := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"infamy"}}, seat("cain")},
		Monsters: []engine.CardRef{"trite", "fly"},
	})
	tb2.Attack("trite", 3, 6)
	if d := tb2.G.Players[0].Damage; d != 1 {
		t.Errorf("damage = %d, want 1 on an attack roll of 3", d)
	}
}
