package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestStoney(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac"), seat("cain")}, Monsters: items("fly", "stoney")}, Set)
	stoney, _ := tb.G.Monsters[1].TopOf()
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	if opts := tb.Options(); slices.Contains(opts, "stoney") {
		t.Errorf("Stoney can be attacked: %v", opts)
	}
	tb.G.ForceRolls(6, 6)
	tb.Settle("fly")
	if tb.G.Object(stoney).Zone.Kind == engine.ZoneInPlay {
		t.Error("Stoney did not die with the fly")
	}
}
