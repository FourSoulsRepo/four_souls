package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestXXJudgement(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"xx_judgement"}},
		{Character: "cain", Souls: []engine.CardRef{"monstro", "gurdy"}},
	}}, Set)
	tb.Play(0, "xx_judgement", foe, "gurdy")
	if s := tb.G.SoulValue(1); s != 1 {
		t.Errorf("soul value = %d, want 1", s)
	}
}
