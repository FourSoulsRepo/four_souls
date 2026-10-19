package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCharmedFatty(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"bomb"}}, seat("cain")},
		Monsters: []engine.CardRef{"charmed_fatty"},
	})
	hurt(tb, 0, 1)
	tb.G.ForceRolls(1, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "charmed_fatty", "bomb", foe)
	if handHas(tb.G, 0, "bomb") {
		t.Error("its damage made Isaac discard the bomb")
	}
}
