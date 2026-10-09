package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHostHat(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac", "host_hat"), seat("cain"))
	tb.Activate(0, "host_hat", 0)
	tb.G.ForceRolls(1, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", foe)
	if a, b := tb.G.Players[0].Damage, tb.G.Players[1].Damage; a != 0 || b != 1 {
		t.Errorf("damage %d and %d, want 0 and 1", a, b)
	}
}
