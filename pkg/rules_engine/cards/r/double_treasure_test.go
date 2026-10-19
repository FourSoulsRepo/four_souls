package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDoubleTreasure(t *testing.T) {
	tb := fight(t, "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "double_treasure")
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", foe)
	if len(tb.G.Players[0].InPlay) != 1 || len(tb.G.Players[1].InPlay) != 1 {
		t.Error("Isaac and Cain each gain +1 treasure")
	}
}
