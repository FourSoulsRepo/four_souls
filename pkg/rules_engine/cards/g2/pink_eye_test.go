package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPinkEye(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"pink_eye"}}, seat("cain")},
		Monsters: []engine.CardRef{"leech", "fly"},
	})
	// A miss on a 1: the fly deals 1 combat damage, Pink Eye deals 1 back
	// and kills it (1 HP).
	tb.Attack("fly", 1)
	if top, _ := tb.G.Monsters[1].TopOf(); tb.G.Object(top).Card == "fly" {
		t.Error("Pink Eye did not damage the fly")
	}
	if d := tb.G.Players[0].Damage; d != 1 {
		t.Errorf("Isaac's damage = %d, want 1", d)
	}
}
