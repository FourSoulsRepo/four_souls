package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestIIITheEmpress(t *testing.T) {
	tb := lootTable(t, "iii_the_empress")
	tb.Play(0, "iii_the_empress", me)
	if got := tb.G.PlayerATK(0); got != 2 {
		t.Errorf("ATK = %d, want 2", got)
	}
	leech, _ := tb.G.Monsters[1].TopOf()
	tb.Attack("leech", 3) // 3 + 1 hits the leech's DC 4; 2 ATK kills it
	if tb.G.Object(leech).Zone.Kind == engine.ZoneInPlay {
		t.Error("the +1 to dice rolls did not apply to the attack roll")
	}
}
