package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSacredHeart(t *testing.T) {
	tb := itemTable(t, items("sacred_heart"))
	fly, _ := tb.G.Monsters[0].TopOf()
	tb.G.ForceRolls(1)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", "yes")
	if tb.G.Object(fly).Zone.Kind == engine.ZoneInPlay {
		t.Error("the 1 did not become a 6")
	}
}
