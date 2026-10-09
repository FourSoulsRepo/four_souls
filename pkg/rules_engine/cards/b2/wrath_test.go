package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestWrath(t *testing.T) {
	tb := slayTable(t, "wrath", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6, 6, 6, 1)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "wrath")
	if a, b := tb.G.Players[0].Damage, tb.G.Players[1].Damage; a != 1 || b != 1 {
		t.Errorf("damage %d and %d, want 1 each", a, b)
	}
}
