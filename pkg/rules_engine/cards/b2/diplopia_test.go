package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestDiplopia(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "diplopia"), seat("cain", "breakfast", "jawbone")}}, Set)
	tb.Start(0, "diplopia", 0)
	if opts := tb.Options(); len(opts) != 2 || opts[0] != "breakfast" {
		t.Fatalf("options %v, want the breakfast (passive) and cancel", opts)
	}
	tb.Settle("breakfast")
	if hp := tb.G.PlayerHP(0); hp != 3 {
		t.Errorf("HP = %d, want 3 (a second breakfast)", hp)
	}
	tb.EndTurn()
	if hp := tb.G.PlayerHP(0); hp != 2 {
		t.Errorf("HP after the turn = %d, want 2", hp)
	}
}
