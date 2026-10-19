package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestIsaacsBlessing(t *testing.T) {
	tb := roomTable(t, "isaacs_blessing")
	tb.G.ForceRolls(1, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "yes", "no")
	if !inDiscard(tb.G, engine.MonsterDeck, "fly") {
		t.Error("Isaac rerolls the 1 and keeps the 6")
	}
}
