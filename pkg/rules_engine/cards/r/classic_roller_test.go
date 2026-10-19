package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestClassicRoller(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("flash_isaac", "classic_roller"), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"leech"},
	})
	tb.G.ForceRolls(1, 6, 2)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	tb.Pass(2)
	tb.Choose("leech")
	tb.ThrowPending() // a 1: the attack roll waits
	tb.Activate(0, "classic_roller", 0, "roll of 1", "player 3 (samson)", "6")
	if top(tb, 0) == "leech" {
		t.Error("Samson picks the 6: a hit kills the leech")
	}
}
