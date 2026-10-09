package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestXWheelOfFortune(t *testing.T) {
	for roll, check := range map[int]func(g *engine.Game) bool{
		1: func(g *engine.Game) bool { return g.Players[0].Cents == 5 },
		2: func(g *engine.Game) bool { return g.Players[0].Dead },
		3: func(g *engine.Game) bool { return len(g.Players[0].Hand) == 3 },
		4: func(g *engine.Game) bool { return g.Players[0].Cents == 0 },
		5: func(g *engine.Game) bool { return g.Players[0].Cents == 9 },
		6: func(g *engine.Game) bool { return len(g.Players[0].InPlay) == 1 },
	} {
		tb := lootTable(t, "x_wheel_of_fortune")
		tb.G.Players[0].Cents = 4
		tb.G.ForceRolls(roll)
		tb.Play(0, "x_wheel_of_fortune")
		if !check(tb.G) {
			t.Errorf("roll %d: wrong result", roll)
		}
	}
}
