package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestBellyButton(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: items("belly_button"), Deactivated: items("character")}, seat("cain")},
		Monsters: []engine.CardRef{"fly"},
	}, Set)
	tb.G.ForceRolls(1, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "yes")
	if !tb.G.Object(tb.G.Players[0].Character).Charged {
		t.Error("the character was not recharged")
	}
}
