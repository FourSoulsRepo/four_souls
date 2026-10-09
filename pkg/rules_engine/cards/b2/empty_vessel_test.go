package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestEmptyVessel(t *testing.T) {
	tb := itemTable(t, items("empty_vessel"))
	if a := tb.G.PlayerATK(0); a != 2 {
		t.Errorf("ATK with no cards = %d, want 2", a)
	}
	leech, _ := tb.G.Monsters[1].TopOf()
	tb.Attack("leech", 3) // 0¢: +1 to the attack roll hits DC 4
	if tb.G.Object(leech).Zone.Kind == engine.ZoneInPlay {
		t.Error("+1 to attack rolls with 0¢ did not apply")
	}
}
