package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCain(t *testing.T) {
	testExtraLoot(t, "cain")
	for _, chars := range [][]engine.CardRef{{"isaac", "cain"}, {"cain", "eve", "maggy"}} {
		g := newGame(t, chars...)
		want := engine.PlayerID(0)
		for i, c := range chars {
			if c == "cain" {
				want = engine.PlayerID(i)
			}
		}
		if g.Turn.Active != want {
			t.Errorf("%v: player %d goes first, want Cain's player %d", chars, g.Turn.Active, want)
		}
		if !hasItem(g, int(want), "sleight_of_hand") {
			t.Error("Cain has no Sleight of Hand")
		}
	}
}
