package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestImp(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"placenta"}}, seat("cain")},
		Monsters: []engine.CardRef{"imp", "fly"},
	})
	tb.G.ForceRolls(3, 5, 6, 6, 6) // a 3: DC 5 → 6, so the 5 misses
	ev := tb.Do(attackIntent(tb), "imp")
	if n := attackRolls(ev); n != 5 {
		t.Errorf("attack rolls = %d, want 5: after a 3, the 5 misses", n)
	}
	if d := tb.G.Players[0].Damage; d != 2 {
		t.Errorf("damage = %d, want 2 from the two misses", d)
	}
}
