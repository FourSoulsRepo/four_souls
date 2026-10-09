package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDaddyHaunt(t *testing.T) {
	tb := monsterTable(t, items("fly"), seat("isaac", "daddy_haunt"), seat("cain"))
	tb.G.ForceRolls(1) // a miss: 1 + 1 damage kills Isaac; the haunt goes to Cain
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", foe)
	if !tb.G.Players[0].Dead || !hasItem(tb.G, 1, "daddy_haunt") {
		t.Errorf("dead %v, haunt with Cain %v", tb.G.Players[0].Dead, hasItem(tb.G, 1, "daddy_haunt"))
	}
}
