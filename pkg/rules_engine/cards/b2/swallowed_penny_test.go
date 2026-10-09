package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSwallowedPenny(t *testing.T) {
	tb := monsterTable(t, []engine.CardRef{"fly"}, seat("isaac", "swallowed_penny"), seat("cain"))
	tb.Attack("fly", 1, 6)
	if c := tb.G.Players[0].Cents; c != 2 { // 1 for the damage, 1 fly reward
		t.Errorf("cents = %d, want 2", c)
	}
}
