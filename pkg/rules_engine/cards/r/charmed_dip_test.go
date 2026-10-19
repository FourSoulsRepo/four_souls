package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCharmedDip(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), {Character: "cain", Items: []engine.CardRef{"jawbone"}, Deactivated: []engine.CardRef{"jawbone"}}},
		Monsters: []engine.CardRef{"charmed_dip"},
	})
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "charmed_dip", foe, "jawbone")
	if !tb.G.Object(tb.Find(1, "jawbone")).Charged {
		t.Error("Cain recharges his Jawbone")
	}
}
