package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDarkOne(t *testing.T) {
	tb := slayTable(t, "dark_one", seat("isaac"), seat("cain"))
	tb.G.ForceRolls(6, 1)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "dark_one")
	if !tb.G.Players[0].Dead && tb.G.Players[0].Damage < 2 {
		t.Errorf("damage = %d: +1 ATK after taking damage did not apply", tb.G.Players[0].Damage)
	}
}
