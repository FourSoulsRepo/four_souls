package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLodestone(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "lodestone"), seat("cain")}, Monsters: []engine.CardRef{"pooter"}})
	tb.G.ForceRolls(3)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "pooter", "pooter")
	// The 3 missed (it resolved before Lodestone); the next roll needs 2+.
	if dc := tb.G.Evasion(tb.G.Monsters[0].Cards[len(tb.G.Monsters[0].Cards)-1]); dc != 2 && top(tb, 0) == "pooter" {
		t.Errorf("the pooter has DC %d, want 3 - 1", dc)
	}
}
