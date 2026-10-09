package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestTheHabit(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: items("the_habit", "jawbone"), Deactivated: items("jawbone")}, seat("cain")},
		Monsters: items("fly"),
	}, Set)
	tb.G.ForceRolls(1, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "jawbone")
	if !tb.G.Object(tb.Find(0, "jawbone")).Charged {
		t.Error("the jawbone was not recharged")
	}
}
