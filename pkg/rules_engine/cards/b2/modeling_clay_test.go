package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestModelingClay(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "modeling_clay"), seat("cain", "breakfast")}}, Set)
	tb.Activate(0, "modeling_clay", 0, "breakfast")
	tb.EndTurn()
	if hp := tb.G.PlayerHP(0); hp != 3 {
		t.Errorf("HP = %d, want 3: the copy stays", hp)
	}
}
