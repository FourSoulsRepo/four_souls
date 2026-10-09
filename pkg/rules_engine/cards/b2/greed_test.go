package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGreed(t *testing.T) {
	tb := slayTable(t, "greed", engine.SituationPlayer{Character: "isaac", Items: items("breakfast"), Cents: 5}, engine.SituationPlayer{Character: "cain", Cents: 3})
	tb.G.ForceRolls(1, 6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "greed")
	if a, b := tb.G.Players[0].Cents, tb.G.Players[1].Cents; a != 1+9 || b != 0 {
		t.Errorf("cents %d and %d, want 10 (1 + the 9¢ reward) and 0", a, b)
	}
}
