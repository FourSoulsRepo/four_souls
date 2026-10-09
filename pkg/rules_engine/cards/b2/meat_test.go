package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMeat(t *testing.T) {
	tb := itemTable(t, items("meat"))
	leech, _ := tb.G.Monsters[1].TopOf()
	tb.Attack("leech", 3) // 3 + 1 hits DC 4
	if tb.G.Object(leech).Zone.Kind == engine.ZoneInPlay {
		t.Error("+1 to attack rolls did not apply")
	}
}
